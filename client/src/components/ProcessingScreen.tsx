import React, { useEffect } from 'react';
import { useNavigate } from 'react-router-dom';

interface ProcessingScreenProps {
  setError: (error: string | null) => void;
}

const ProcessingScreen: React.FC<ProcessingScreenProps> = ({ setError }) => {
  const navigate = useNavigate();

  useEffect(() => {
    if (localStorage.getItem('initialText')) {
      navigate('/chat')
    }
    
    const timeoutId = setTimeout(() => {
      setError('Request timed out. Please try again.');
      navigate('/');
    }, 30000); // 30 seconds timeout

    return () => clearTimeout(timeoutId);
  }, [navigate, setError]);

  return (
    <div className="flex items-center justify-center min-h-screen">
      <div className="text-center">
        <div className="animate-spin rounded-full h-32 w-32 border-t-2 border-b-2 border-dark-primary mx-auto mb-4"></div>
        <div className="text-2xl font-bold text-dark-primary">Processing...</div>
      </div>
    </div>
  );
};

export default ProcessingScreen;