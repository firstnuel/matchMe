import React from 'react'

interface IsLoadingProps {
  message?: string
  size?: 'small' | 'medium' | 'large'
}

const IsLoading: React.FC<IsLoadingProps> = ({ 
  message = 'Loading...', 
  size = 'medium' 
}) => {
  return (
    <div className={`loading-container loading-${size}`}>
      <div className="loading-spinner"></div>
      <span className="loading-message">{message}</span>
    </div>
  )
}

export default IsLoading