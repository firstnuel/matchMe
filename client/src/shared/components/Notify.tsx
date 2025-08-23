import { useEffect } from 'react';
import { useUIStore } from '../hooks/uiStore';
import '../../App.css';

const Notify = () => {
  const { errorMsg, infoMsg, clearMsgs } = useUIStore();
  const message = errorMsg || infoMsg;
  const type = errorMsg ? 'error' : infoMsg ? 'success' : null;

  useEffect(() => {
    if (message) {
      console.log('Setting timeout for message:', message); // Debug log
      const timer = setTimeout(() => {
        console.log('Timeout triggered, clearing messages'); // Debug log
        clearMsgs();
      }, 3000);
      return () => {
        console.log('Cleaning up timeout'); // Debug log
        clearTimeout(timer);
      };
    }
  }, [message, clearMsgs]);

  if (!message || !type) {
    console.log('Not rendering Notify, message or type is null'); // Debug log
    return null;
  }

  console.log('Rendering Notify with message:', message, 'type:', type); // Debug log

  return (
    <div className={`notify-${type}-div`}>
      <div className={`${type}-msg`}>{message}</div>
      <div className={`clear-${type}-btn`} onClick={clearMsgs}>
        &times;
      </div>
    </div>
  );
};

export default Notify;