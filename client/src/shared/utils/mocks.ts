import type { Chat, ChatUser } from "../../features/chat/types/chat";

// Mock data
export const mockCurrentUser: ChatUser = {
  id: '1',
  firstName: 'John',
  lastName: 'Doe',
  isOnline: true
};

export const mockChats: Chat[] = [
  {
    id: '1',
    participants: [
      mockCurrentUser,
      { id: '2', firstName: 'Sarah', lastName: 'M.', isOnline: true }
    ],
    unreadCount: 2,
    updatedAt: new Date(Date.now() - 2 * 60 * 1000), // 2 minutes ago
    lastMessage: {
      id: '1',
      senderId: '2',
      receiverId: '1',
      content: 'Hey! Thanks for the match 😊',
      timestamp: new Date(Date.now() - 2 * 60 * 1000),
      status: 'delivered',
      type: 'text'
    },
    messages: [
      {
        id: '1',
        senderId: '2',
        receiverId: '1',
        content: 'Hey John! 👋 Thanks for the match! I love your hiking photos',
        timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000),
        status: 'read',
        type: 'text'
      },
      {
        id: '2',
        senderId: '1',
        receiverId: '2',
        content: 'Hey Sarah! Thank you 😊 I noticed you\'re into photography too. Your sunset shots are incredible!',
        timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000 + 3 * 60 * 1000),
        status: 'read',
        type: 'text'
      },
      {
        id: '3',
        senderId: '2',
        receiverId: '1',
        content: 'Aww thanks! 📸 I\'m always chasing the perfect golden hour. Do you have a favorite hiking spot?',
        timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000 + 5 * 60 * 1000),
        status: 'read',
        type: 'text'
      },
      {
        id: '4',
        senderId: '1',
        receiverId: '2',
        content: 'Definitely! There\'s this amazing trail in the Blue Ridge Mountains. The view at the summit is breathtaking, especially during sunrise 🌅',
        timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000 + 8 * 60 * 1000),
        status: 'read',
        type: 'text'
      },
      {
        id: '5',
        senderId: '2',
        receiverId: '1',
        content: 'That sounds perfect! I\'ve been wanting to explore more mountain trails. Maybe we could plan a hiking adventure sometime? ⛰️',
        timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000 + 10 * 60 * 1000),
        status: 'read',
        type: 'text'
      },
      {
        id: '6',
        senderId: '1',
        receiverId: '2',
        content: 'I\'d love that! How about this weekend? Weather\'s supposed to be perfect ☀️',
        timestamp: new Date(Date.now() - 2 * 60 * 60 * 1000 + 12 * 60 * 1000),
        status: 'delivered',
        type: 'text'
      }
    ]
  },
  {
    id: '2',
    participants: [
      mockCurrentUser,
      { id: '3', firstName: 'Emma', lastName: 'J.', isOnline: true }
    ],
    unreadCount: 0,
    updatedAt: new Date(Date.now() - 60 * 60 * 1000), // 1 hour ago
    lastMessage: {
      id: '2',
      senderId: '3',
      receiverId: '1',
      content: 'That coffee place you mentioned sounds amazing!',
      timestamp: new Date(Date.now() - 60 * 60 * 1000),
      status: 'read',
      type: 'text'
    },
    messages: []
  },
  {
    id: '3',
    participants: [
      mockCurrentUser,
      { id: '4', firstName: 'Alex', lastName: 'L.', isOnline: false, lastSeen: new Date(Date.now() - 2 * 60 * 60 * 1000) }
    ],
    unreadCount: 0,
    updatedAt: new Date(Date.now() - 3 * 60 * 60 * 1000), // 3 hours ago
    lastMessage: {
      id: '3',
      senderId: '4',
      receiverId: '1',
      content: 'Would love to hear more about your hiking adventures',
      timestamp: new Date(Date.now() - 3 * 60 * 60 * 1000),
      status: 'read',
      type: 'text'
    },
    messages: []
  }
];


