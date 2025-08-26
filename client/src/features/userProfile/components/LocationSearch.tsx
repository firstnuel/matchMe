import { useState, useEffect, useRef } from "react";
import { Icon } from "@iconify/react/dist/iconify.js";
import { type Point } from "../../../shared/types/user";

interface LocationSearchProps {
  location: Point;
  setLocation: React.Dispatch<React.SetStateAction<Point>>;
}

interface LocationResult {
  display_name: string;
  lat: string;
  lon: string;
  name: string;
  type: string;
}

const LocationSearch = ({ location, setLocation }: LocationSearchProps) => {
  const [searchQuery, setSearchQuery] = useState("");
  const [suggestions, setSuggestions] = useState<LocationResult[]>([]);
  const [isLoading, setIsLoading] = useState(false);
  const [showDropdown, setShowDropdown] = useState(false);
  const [error, setError] = useState("");
  const timeoutRef = useRef<ReturnType<typeof setTimeout> | null>(null);
  const dropdownRef = useRef<HTMLDivElement>(null);

  const searchLocations = async (query: string) => {
    if (!query.trim()) {
      setSuggestions([]);
      return;
    }

    setIsLoading(true);
    setError("");

    try {
      const response = await fetch(
        `https://nominatim.openstreetmap.org/search?format=json&q=${encodeURIComponent(
          query
        )}&limit=5&addressdetails=1`
      );
      
      if (!response.ok) {
        throw new Error("Failed to search locations");
      }
      
      const data: LocationResult[] = await response.json();
      setSuggestions(data);
      setShowDropdown(true);
    } catch (err) {
      console.error("Location search error:", err);
      setError("Failed to search locations. Please try again.");
      setSuggestions([]);
    } finally {
      setIsLoading(false);
    }
  };

  const handleInputChange = (value: string) => {
    setSearchQuery(value);
    
    if (timeoutRef.current) {
      clearTimeout(timeoutRef.current);
    }

    timeoutRef.current = setTimeout(() => {
      searchLocations(value);
    }, 300);
  };

  const handleLocationSelect = (locationResult: LocationResult) => {
    const cityName = locationResult.display_name.split(",")[0] || locationResult.name;
    
    setLocation({
      latitude: parseFloat(locationResult.lat),
      longitude: parseFloat(locationResult.lon),
      city: cityName,
    });
    
    setSearchQuery(cityName);
    setShowDropdown(false);
    setSuggestions([]);
    setError("");
  };

  const handleClickOutside = (event: MouseEvent) => {
    if (dropdownRef.current && !dropdownRef.current.contains(event.target as Node)) {
      setShowDropdown(false);
    }
  };

  useEffect(() => {
    document.addEventListener("mousedown", handleClickOutside);
    return () => {
      document.removeEventListener("mousedown", handleClickOutside);
      if (timeoutRef.current) {
        clearTimeout(timeoutRef.current);
      }
    };
  }, []);

  useEffect(() => {
    if (location.city && !searchQuery) {
      setSearchQuery(location.city);
    }
  }, [location.city, searchQuery]);

  return (
    <div className="location-section">
      <div>
        <Icon icon="mdi:location" className="location-icon" />
        Search for your location
      </div>
      
      <div className="location-search-container" ref={dropdownRef} style={{ position: "relative" }}>
        <input
          type="text"
          className="location-search-input"
          placeholder="Search for a city or location..."
          value={searchQuery}
          onChange={(e) => handleInputChange(e.target.value)}
          onFocus={() => {
            if (suggestions.length > 0) {
              setShowDropdown(true);
            }
          }}
          style={{
            width: "100%",
            padding: "10px 40px 10px 12px",
            border: "1px solid #ddd",
            borderRadius: "6px",
            fontSize: "14px",
            outline: "none",
          }}
        />
        
        {isLoading && (
          <Icon 
            icon="mdi:loading" 
            className="location-search-loading"
            style={{
              position: "absolute",
              right: "12px",
              top: "50%",
              transform: "translateY(-50%)",
              animation: "spin 1s linear infinite"
            }}
          />
        )}
        
        {showDropdown && suggestions.length > 0 && (
          <div 
            className="location-dropdown"
            style={{
              position: "absolute",
              top: "100%",
              left: 0,
              right: 0,
              backgroundColor: "white",
              border: "1px solid #ddd",
              borderTop: "none",
              borderRadius: "0 0 6px 6px",
              maxHeight: "200px",
              overflowY: "auto",
              zIndex: 1000,
              boxShadow: "0 2px 8px rgba(0,0,0,0.1)"
            }}
          >
            {suggestions.map((suggestion, index) => (
              <div
                key={index}
                className="location-dropdown-item"
                onClick={() => handleLocationSelect(suggestion)}
                style={{
                  padding: "12px",
                  cursor: "pointer",
                  borderBottom: index < suggestions.length - 1 ? "1px solid #f0f0f0" : "none",
                  fontSize: "14px",
                }}
                onMouseEnter={(e) => {
                  (e.target as HTMLDivElement).style.backgroundColor = "#f5f5f5";
                }}
                onMouseLeave={(e) => {
                  (e.target as HTMLDivElement).style.backgroundColor = "white";
                }}
              >
                <div style={{ fontWeight: "500" }}>
                  {suggestion.display_name.split(",")[0]}
                </div>
                <div style={{ fontSize: "12px", color: "#666", marginTop: "2px" }}>
                  {suggestion.display_name}
                </div>
              </div>
            ))}
          </div>
        )}
      </div>
      
      <div className="current-location">
        Current: {location.city || "No location selected"}{" "}
        {location.latitude && location.longitude
          ? `(${location.latitude.toFixed(4)}, ${location.longitude.toFixed(4)})`
          : ""}
      </div>
      
      {error && <div className="info error">{error}</div>}
      
      <style>{`
        @keyframes spin {
          from { transform: translateY(-50%) rotate(0deg); }
          to { transform: translateY(-50%) rotate(360deg); }
        }
      `}</style>
    </div>
  );
};

export default LocationSearch;