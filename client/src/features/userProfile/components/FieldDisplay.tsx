interface FieldDisplayProps {
  label: string;
  value: string | number | string[] | null | undefined;
  type?: 'text' | 'array' | 'range';
  className?: string
}

const FieldDisplay = ({ label, value, className="", type = 'text' }: FieldDisplayProps) => {
  const formatValue = () => {
    if (value === null || value === undefined || value === '') {
      return 'Not specified';
    }

    if (type === 'array' && Array.isArray(value)) {
      return value.length > 0 ? value.join(', ') : 'Not specified';
    }

    if (type === 'range' && Array.isArray(value) && value.length === 2) {
      return `${value[0]} - ${value[1]}`;
    }

    return String(value);
  };

  return (
    <div className={`sec-show ${className}`}>
      <div className="sec-name">{label}</div>
      <div className="sec-value">{formatValue()}</div>
    </div>
  );
};

export default FieldDisplay;