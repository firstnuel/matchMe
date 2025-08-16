package schema

import (
	"bytes"
	"database/sql/driver"
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"strings"

	"entgo.io/ent"
	"entgo.io/ent/dialect"
)

// Point represents a PostGIS geography point (lon, lat).
type Point struct {
	Longitude float64
	Latitude  float64
}

// Scan implements the Scanner interface for Point.
func (p *Point) Scan(src interface{}) error {
	if src == nil {
		return nil
	}

	var s string
	// Handle multiple source types from the database driver.
	// Some drivers might return a hex string as []byte.
	switch v := src.(type) {
	case string:
		s = v
	case []byte:
		s = string(v)
	default:
		return fmt.Errorf("unsupported scan type for Point: %T", src)
	}

	s = strings.TrimSpace(s)

	// Try WKB hex format first (common from PostGIS)
	if strings.HasPrefix(s, "01") || strings.HasPrefix(s, "00") {
		return p.scanWKB(s)
	}

	// Fallback to parsing WKT format: "POINT(longitude latitude)"
	s = strings.TrimPrefix(s, "SRID=4326;")
	var lon, lat float64
	_, err := fmt.Sscanf(s, "POINT(%f %f)", &lon, &lat)
	if err != nil {
		return fmt.Errorf("failed to parse WKT point: %w", err)
	}
	p.Longitude = lon
	p.Latitude = lat
	return nil
}

// scanWKB parses WKB hex string
func (p *Point) scanWKB(s string) error {
	b, err := hex.DecodeString(s)
	if err != nil {
		return fmt.Errorf("failed to decode WKB hex: %w", err)
	}
	return p.scanWKBBytes(b)
}

// scanWKBBytes parses WKB binary format
func (p *Point) scanWKBBytes(b []byte) error {
	if len(b) < 21 { // Minimum size for a 2D WKB point
		return fmt.Errorf("invalid WKB point: too short, got %d bytes", len(b))
	}

	r := bytes.NewReader(b)

	// Read byte order (1 byte)
	var byteOrder byte
	if err := binary.Read(r, binary.LittleEndian, &byteOrder); err != nil {
		return fmt.Errorf("failed to read WKB byte order: %w", err)
	}

	var order binary.ByteOrder
	if byteOrder == 1 {
		order = binary.LittleEndian // NDR (little endian)
	} else {
		order = binary.BigEndian // XDR (big endian)
	}

	// Read the full geometry type word (4 bytes)
	var typeWord uint32
	if err := binary.Read(r, order, &typeWord); err != nil {
		return fmt.Errorf("failed to read WKB geometry type word: %w", err)
	}

	// In EWKB, the type word contains flags. Check for and skip the SRID if the flag is set.
	// The SRID flag is 0x20000000
	if typeWord&0x20000000 != 0 {
		var srid uint32
		if err := binary.Read(r, order, &srid); err != nil {
			return fmt.Errorf("failed to read EWKB SRID: %w", err)
		}
	}

	// Mask out all flags to get the geometry type, including Z/M dimension info
	// Use 0x1FFFFFFF to correctly mask out the high-bit flags (Z, M, SRID)
	geomType := typeWord & 0x1FFFFFFF

	// Check if it's a Point type (2D, Z, M, or ZM)
	// WKB Point types are 1 (Point), 1001 (PointZ), 2001 (PointM), 3001 (PointZM)
	// A simple modulo check handles all cases.
	if geomType%1000 != 1 {
		return fmt.Errorf("invalid WKB type: expected a Point, but got type %d (from raw type word %d)", geomType, typeWord)
	}

	// Read coordinates (8 bytes each)
	if err := binary.Read(r, order, &p.Longitude); err != nil {
		return fmt.Errorf("failed to read longitude: %w", err)
	}
	if err := binary.Read(r, order, &p.Latitude); err != nil {
		return fmt.Errorf("failed to read latitude: %w", err)
	}

	// Note: This implementation assumes a 2D point and will ignore Z/M coordinates
	// if they are present in the data. To support Z/M, you would need to check
	// the geomType and read additional float64 values from the reader.

	return nil
}

// Value implements the Valuer interface for Point.
func (p Point) Value() (driver.Value, error) {
	if p.Longitude == 0 && p.Latitude == 0 {
		return nil, nil
	}
	return fmt.Sprintf("SRID=4326;POINT(%f %f)", p.Longitude, p.Latitude), nil
}

// String returns a string representation of the Point.
func (p Point) String() string {
	return fmt.Sprintf("POINT(%f %f)", p.Longitude, p.Latitude)
}

// PostGISDialect is a custom dialect extension for PostGIS.
type PostGISDialect struct {
	dialect.Driver
}

// NewPostGISDialect creates a new PostGIS dialect wrapping the given driver.
func NewPostGISDialect(driver dialect.Driver) *PostGISDialect {
	return &PostGISDialect{Driver: driver}
}

// SchemaType defines the SQL type for PostGIS fields.
func (d *PostGISDialect) SchemaType(f ent.Field) (string, error) {
	// Check for custom schema type annotation
	if schemaType, ok := f.Descriptor().SchemaType[dialect.Postgres]; ok {
		if strings.Contains(schemaType, "geography") || strings.Contains(schemaType, "geometry") {
			return schemaType, nil
		}
	}

	// Delegate to wrapped driver for non-PostGIS types
	if d.Driver != nil {
		switch drv := d.Driver.(type) {
		case interface {
			SchemaType(ent.Field) (string, error)
		}:
			return drv.SchemaType(f)
		}
	}

	return "", fmt.Errorf("unsupported field type: %v", f.Descriptor().Info)
}

// PostGISExtension represents the PostGIS extension migration.
type PostGISExtension struct{}

// Name returns the migration name.
func (PostGISExtension) Name() string {
	return "create_postgis_extension"
}

// SQL returns the SQL statement to create the PostGIS extension.
func (PostGISExtension) SQL() string {
	return "CREATE EXTENSION IF NOT EXISTS postgis"
}
