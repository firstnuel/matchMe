import React, { useState } from 'react';
import ConnectionItem from './ConnectionItem';
import { useConnections, useConnectionRequests } from '../hooks/useConnections';
import { getInitials } from '../../../shared/utils/utils';
import '../styles.css'

type TabType = 'requests' | 'matches';

const Connection: React.FC = () => {
  const [activeTab, setActiveTab] = useState<TabType>('requests');
  
  const { data: connectionsData, isLoading: connectionsLoading } = useConnections();
  const { data: requestsData, isLoading: requestsLoading } = useConnectionRequests();
  
  const connections = connectionsData && 'connections' in connectionsData ? connectionsData.connections : [];
  const requests = requestsData && 'requests' in requestsData ? requestsData.requests : [];
  
  const handleTabChange = (tab: TabType) => {
    setActiveTab(tab);
  };

  const renderContent = () => {
    if (activeTab === 'requests') {
      if (requestsLoading) {
        return <div className="loading-state">Loading requests...</div>;
      }
      
      if (requests.length === 0) {
        return (
          <div className="empty-state">
            <h3>No Connection Requests</h3>
            <p>You don't have any pending connection requests at the moment.</p>
          </div>
        );
      }
      
      return requests.map((request) => (
        <ConnectionItem
          key={request.id}
          id={request.id}
          type="request"
          user={request.sender}
          initials={getInitials(request.sender?.first_name || '', request.sender?.last_name || '')}
          name={`${request.sender?.first_name || 'Unknown'} ${request.sender?.last_name?.[0] || ''}.`}
          description={request.message || "Wants to connect"}
          createdAt={request.created_at}
        />
      ));
    }
    
    // Matches tab
    if (connectionsLoading) {
      return <div className="loading-state">Loading connections...</div>;
    }
    
    if (connections.length === 0) {
      return (
        <div className="empty-state">
          <h3>No Connections Yet</h3>
          <p>Start connecting with people to see your matches here!</p>
        </div>
      );
    }
    
    return connections.map((connection) => {
      // Determine which user to show (not the current user)
      const otherUser = connection.user_a || connection.user_b;
      
      return (
        <ConnectionItem
          key={connection.id}
          id={connection.id}
          type="connection"
          user={otherUser}
          initials={getInitials(otherUser?.first_name || '', otherUser?.last_name || '')}
          name={`${otherUser?.first_name || 'Unknown'} ${otherUser?.last_name?.[0] || ''}.`}
          description={`Connected on ${new Date(connection.connected_at).toLocaleDateString()}`}
          createdAt={connection.connected_at}
        />
      );
    });
  };

  return (
    <div id="connections" className="view">
      <div className="connections-view">
        <div className="connection-tabs">
          <button 
            className={`tab-btn ${activeTab === 'requests' ? 'active' : ''}`}
            onClick={() => handleTabChange('requests')}
          >
            Requests ({requests.length})
          </button>
          <button 
            className={`tab-btn ${activeTab === 'matches' ? 'active' : ''}`}
            onClick={() => handleTabChange('matches')}
          >
            Matches ({connections.length})
          </button>
        </div>
        
        <div className="connections-content">
          {renderContent()}
        </div>
      </div>
    </div>
  );
};

export default Connection;