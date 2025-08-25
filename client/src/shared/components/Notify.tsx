import { useEffect } from 'react';
import { useUIStore } from '../hooks/uiStore';
import '../../App.css';

const Notify = () => {
  const { errorMsg, infoMsg, clearMsgs } = useUIStore();
  const message = errorMsg || infoMsg;
  const type = errorMsg ? 'error' : infoMsg ? 'success' : null;

  useEffect(() => {
    if (message) {
      const timer = setTimeout(() => {
        clearMsgs();
      }, 3000);
      return () => {
        clearTimeout(timer);
      };
    }
  }, [message, clearMsgs]);

  if (!message || !type) {
    return null;
  }

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