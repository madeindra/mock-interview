import React, { useState, useEffect } from 'react';

interface NavbarProps {
  showBackIcon?: boolean;
  showForwardIcon?: boolean;
  showStartOver?: boolean;
  onBack?: () => void;
  onForward?: () => void;
  onStartOver?: () => void;
  disableBack?: boolean;
  disableForward?: boolean;
}

interface StatusResponse {
  message: string;
  data: {
    backend: boolean;
    api: boolean | null;
    apiStatus: string | null;
    key: boolean;
  };
}

const Navbar: React.FC<NavbarProps> = ({
  showBackIcon = false,
  showForwardIcon = false,
  showStartOver = false,
  onBack,
  onForward,
  onStartOver,
  disableBack = false,
  disableForward = false
}) => {
  const [status, setStatus] = useState<StatusResponse['data'] | null>(null);
  const [showTooltip, setShowTooltip] = useState(false);

  useEffect(() => {
    const fetchStatus = async () => {
      try {
        const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/chat/status`);
        if (!response.ok) {
          throw new Error('Network response was not ok');
        }
        const data: StatusResponse = await response.json();
        setStatus(data.data);
      } catch (error) {
        console.error('Error fetching status:', error);
        setStatus(null);
      }
    };

    fetchStatus();
    const intervalId = setInterval(fetchStatus, 30000); // Fetch every 30 seconds

    return () => clearInterval(intervalId);
  }, []);

  const getStatusColor = () => {
    if (!status) return 'bg-red-500';
    if (status.backend && status.api === true && status.key) return 'bg-green-500';
    if (status.api === false) return 'bg-orange-500';
    return 'bg-red-500';
  };

  const capitalizeFirstLetter = (string: string) => {
    return string.charAt(0).toUpperCase() + string.slice(1);
  };

  if (!showBackIcon && !showForwardIcon && !showStartOver && !status) {
    return null;
  }

  const handleBack = () => {
    if (!disableBack && onBack) {
      onBack();
    }
  }

  const handleForward = () => {
    if (!disableForward && onForward) {
      onForward();
    }
  }

  const handleStartOver = () => {
    if (onStartOver) {
      const isConfirmed = window.confirm("Are you sure you want to start over the interview?");
      if (isConfirmed) {
        onStartOver();
      }
    }
  };

  return (
    <nav className="bg-dark-surface p-4 flex justify-between items-center relative">
      <div className="flex items-center">
        {showBackIcon && (
          <button
            onClick={handleBack}
            className={`mr-4 ${disableBack ? 'text-gray-500 cursor-not-allowed' : 'text-white hover:text-gray-300'}`}
            disabled={disableBack}
          >
            &#8592; {/* Left arrow */}
          </button>
        )}
        {showForwardIcon && (
          <button
            onClick={handleForward}
            className={`${disableForward ? 'text-gray-500 cursor-not-allowed' : 'text-white hover:text-gray-300'}`}
            disabled={disableForward}
          >
            &#8594; {/* Right arrow */}
          </button>
        )}
      </div>
      <div className="flex items-center">
        <div 
          className={`w-4 h-4 rounded-full mr-4 ${getStatusColor()}`}
          onMouseEnter={() => setShowTooltip(true)}
          onMouseLeave={() => setShowTooltip(false)}
        >
          {showTooltip && status && (
            <div className="absolute top-full right-2 mt-2 p-2 bg-white text-black rounded shadow-lg z-10">
              <p>Database: {status.backend ? 'Up' : 'Down'}</p>
              <p>API: {status.api === null ? 'Unknown' : (status.api ? 'Up' : 'Down')}</p>
              <p>Status: {status.apiStatus ? capitalizeFirstLetter(status.apiStatus) : 'Unknown'}</p>
              <p>Authorized: {status.key ? 'Yes' : 'No'}</p>
            </div>
          )}
        </div>
        {showStartOver && (
          <button 
            onClick={handleStartOver} 
            className="text-white hover:text-gray-300"
          >
            &#x2715; {/* Multiplication symbol */}
          </button>
        )}
      </div>
    </nav>
  );
};

export default Navbar;