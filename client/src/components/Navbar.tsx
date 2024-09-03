import React from 'react';

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
  if (!showBackIcon && !showForwardIcon && !showStartOver) {
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
    <nav className="bg-dark-surface p-4 flex justify-between items-center">
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
      {showStartOver && (
        <button 
          onClick={handleStartOver} 
          className="text-white hover:text-gray-300"
        >
          &#x2715; {/* Multiplication symbol */}
        </button>
      )}
    </nav>
  );
};

export default Navbar;