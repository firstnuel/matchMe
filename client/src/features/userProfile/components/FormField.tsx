import React, { useState } from "react";

interface FormFieldProps {
  label: string;
  type?: "text" | "number" | "select" | "textarea" | "range" | "slider";
  value?: string | number | [number, number];
  placeholder?: string;
  options?: { value: string; label: string }[]; // only for select
  onChange?: (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement> | { target: { value: [number, number] } }
  ) => void;
  maxLength?: number;
  min?: number;
  max?: number;
  unit?: string; // for slider display (e.g. "miles", "km")
}

const FormField: React.FC<FormFieldProps> = ({
  label,
  type = "text",
  value = "",
  placeholder,
  options,
  onChange,
  maxLength,
  min,
  max,
  unit,
}) => {
  const [charCount, setCharCount] = useState(
    typeof value === "string" ? value.length : 0
  );


  const handleChange = (
    e: React.ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>
  ) => {
    if (type === "textarea" && maxLength) {
      setCharCount(e.target.value.length);
    }
    if (type === "range" && Array.isArray(value)) {
      // Do nothing, handled by handleRangeChange
      return;
    }
    onChange?.(e);
  };
  const handleRangeChange = (index: 0 | 1, newValue: number) => {
    const current = Array.isArray(value) ? value : [min || 0, max || 0];
    const updated: [number, number] = [...current] as [number, number];
    updated[index] = newValue;
    onChange?.({ target: { value: updated } });
  };

  return (
    <div className="form-group">
      <label className="form-label">{label}</label>

      {type === "select" ? (
        <select className="form-select" value={value as string} onChange={handleChange}>
          {options?.map((opt) => (
            <option key={opt.value} value={opt.value}>
              {opt.label}
            </option>
          ))}
        </select>
      ) : type === "textarea" ? (
        <>
          <textarea
            className="form-textarea"
            value={value as string}
            onChange={handleChange}
            placeholder={placeholder}
            maxLength={maxLength}
          />
          {maxLength && (
            <div className="character-count">
              {charCount}/{maxLength}
            </div>
          )}
        </>
      ) : type === "range" ? (
        <div className="range-group">
          <input
            type="number"
            className="form-input range-input"
            value={Array.isArray(value) ? value[0] : ""}
            min={min}
            max={max}
            onChange={(e) => handleRangeChange(0, Number(e.target.value))}
          />
          <span className="range-separator">to</span>
          <input
            type="number"
            className="form-input range-input"
            value={Array.isArray(value) ? value[1] : ""}
            min={min}
            max={max}
            onChange={(e) => handleRangeChange(1, Number(e.target.value))}
          />
        </div>
      ) : type === "slider" ? (
        <>
          <input
            type="range"
            className="form-input"
            min={min}
            max={max}
            value={value as number}
            onChange={handleChange}
          />
          <div style={{ textAlign: "center", marginTop: "5px", color: "#666" }}>
            {value} {unit}
          </div>
        </>
      ) : (
        <input
          type={type}
          className="form-input"
          value={value as string | number}
          onChange={handleChange}
          maxLength={maxLength}
          min={min}
          max={max}
          placeholder={placeholder}
        />
      )}
    </div>
  );
};

export default FormField;