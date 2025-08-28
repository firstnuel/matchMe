import React from 'react';

interface DateBannerProps {
  dateLabel: string;
}

const DateBanner: React.FC<DateBannerProps> = ({ dateLabel }) => {
  return (
    <div className="message-date">
      <span>{dateLabel}</span>
    </div>
  );
};

export default DateBanner;