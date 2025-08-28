import { type MessageGroup } from "../types/user";

export const firstToUpper = (s: string): string => {
    return s.length <= 1? s : s.charAt(0).toUpperCase() + s.substring(1, s.length)
}

export function getInitials(firstName?: string | null, lastName?: string | null): string {
  const first = firstName?.charAt(0) || '?';
  const last = lastName?.charAt(0) || '?';
  return (first + last).toUpperCase();
}

export const formatTime = (date: Date): string => {
  const now = new Date();
  const diff = now.getTime() - date.getTime();
  const minutes = Math.floor(diff / (1000 * 60));
  const hours = Math.floor(diff / (1000 * 60 * 60));
  const days = Math.floor(diff / (1000 * 60 * 60 * 24));

  if (minutes < 1) return 'now';
  if (minutes < 60) return `${minutes}m`;
  if (hours < 24) return `${hours}h`;
  return `${days}d`;
};

export const formatDateLabel = (date: Date): string => {
  const today = new Date();
  const yesterday = new Date(today);
  yesterday.setDate(yesterday.getDate() - 1);
  
  const messageDate = new Date(date);
  
  // Reset time to compare only dates
  today.setHours(0, 0, 0, 0);
  yesterday.setHours(0, 0, 0, 0);
  messageDate.setHours(0, 0, 0, 0);
  
  if (messageDate.getTime() === today.getTime()) {
    return 'Today';
  } else if (messageDate.getTime() === yesterday.getTime()) {
    return 'Yesterday';
  } else {
    return messageDate.toLocaleDateString('en-US', { 
      weekday: 'long', 
      year: 'numeric', 
      month: 'long', 
      day: 'numeric' 
    });
  }
};

export const groupMessagesByDate = <T extends { created_at: string }>(messages: T[]): MessageGroup[] => {
  const groups = new Map<string, T[]>();
  
  messages.forEach(message => {
    const messageDate = new Date(message.created_at);
    const dateLabel = formatDateLabel(messageDate);
    
    if (!groups.has(dateLabel)) {
      groups.set(dateLabel, []);
    }
    groups.get(dateLabel)!.push(message);
  });
  
  // Convert to array and sort by date (newest first)
  return Array.from(groups.entries())
    .map(([dateLabel, messages]) => ({ dateLabel, messages }))
    .sort((a, b) => {
      // Sort by the first message's date in each group
      const dateA = new Date(a.messages[0].created_at);
      const dateB = new Date(b.messages[0].created_at);
      return dateA.getTime() - dateB.getTime();
    });
};