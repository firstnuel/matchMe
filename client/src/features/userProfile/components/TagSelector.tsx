interface TagSelectorProps {
  options: string[];
  maxSelectable: number;
  label: string;
  className?: string;
  selectedTags: string[];
  setSelectedTags: React.Dispatch<React.SetStateAction<string[]>>
}

const TagSelector: React.FC<TagSelectorProps> = ({
  options,
  maxSelectable,
  label,
  selectedTags,
  setSelectedTags,
  className = ""
}) => {


  const handleTagClick = (tag: string) => {
    if (selectedTags.includes(tag)) {
      // Remove tag if already selected
      setSelectedTags(selectedTags.filter(t => t !== tag));
    } else if (selectedTags.length < maxSelectable) {
      // Add tag if not at max limit
      setSelectedTags([...selectedTags, tag]);
    }
  };

  const isTagSelected = (tag: string) => selectedTags.includes(tag);

  return (
    <div className={`form-group ${className}`}>
      <label className="form-label">{label}</label>
      <div className="tag-container">
        {options.map((option) => (
          <div
            key={option}
            className={`tag ${isTagSelected(option) ? 'tag-selected' : ''}`}
            onClick={() => handleTagClick(option)}
            style={{ cursor: 'pointer' }}
          >
            {option}
            {isTagSelected(option) && (
              <span className="tag-remove" style={{ marginLeft: '5px' }}>
                ×
              </span>
            )}
          </div>
        ))}
      </div>
    </div>
  );
};

export default TagSelector;