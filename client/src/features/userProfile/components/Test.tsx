import { Link } from 'react-router-dom';

const TestRoute = () => {
  return (
    <div>
      <h1>Test Route</h1>
      <p>This is a test route component.</p>
      <Link to="/">Go to Home</Link>
    </div>
  );
};

export default TestRoute;