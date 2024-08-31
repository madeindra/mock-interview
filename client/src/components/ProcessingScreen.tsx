import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

interface ProcessingScreenProps {
  setError: (error: string | null) => void;
}

const ProcessingScreen: React.FC<ProcessingScreenProps> = ({ setError }) => {
  const navigate = useNavigate();

  useEffect(() => {
    const timeoutId = setTimeout(() => {
      setError('Request timed out. Please try again.');
      navigate('/');
    }, 30000); // 30 seconds timeout

    return () => clearTimeout(timeoutId);
  }, [navigate, setError]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-2xl font-bold">Processing...</div>
    </div>
  );
};

export default ProcessingScreen;