import React, { useState } from 'react';
import { BrowserRouter as Router, Routes, Route, Navigate } from 'react-router-dom';
import StartScreen from './components/StartScreen';
import ProcessingScreen from './components/ProcessingScreen';
import ChatScreen from './components/ChatScreen';

const App: React.FC = () => {
  const backendHost = import.meta.env.VITE_BACKEND_URL || 'http://0.0.0.0:8080';
  const [error, setError] = useState<string | null>(null);

  return (
    <Router>
      <div className="min-h-screen bg-dark-bg text-dark-on-bg flex flex-col relative">
        {error && (
          <div className="absolute top-0 left-0 right-0 bg-dark-error text-dark-on-surface px-4 py-3 z-50" role="alert">
            <strong className="font-bold">Error: </strong>
            <span className="block sm:inline">{error}</span>
            <span className="absolute top-0 bottom-0 right-0 px-4 py-3">
              <svg className="fill-current h-6 w-6 text-dark-on-surface" role="button" onClick={() => setError(null)} xmlns="http://www.w3.org/2000/svg" viewBox="0 0 20 20">
                <title>Close</title>
                <path d="M14.348 14.849a1.2 1.2 0 0 1-1.697 0L10 11.819l-2.651 3.029a1.2 1.2 0 1 1-1.697-1.697l2.758-3.15-2.759-3.152a1.2 1.2 0 1 1 1.697-1.697L10 8.183l2.651-3.031a1.2 1.2 0 1 1 1.697 1.697l-2.758 3.152 2.758 3.15a1.2 1.2 0 0 1 0 1.698z"/>
              </svg>
            </span>
          </div>
        )}
        <div className="flex-grow">
          <Routes>
            <Route path="/" element={<StartScreen backendHost={backendHost} setError={setError} />} />
            <Route path="/processing" element={<ProcessingScreen setError={setError} />} />
            <Route path="/chat" element={<ChatScreen  backendHost={backendHost}setError={setError} />} />
            <Route path="*" element={<Navigate to="/" replace />} />
          </Routes>
        </div>
      </div>
    </Router>
  );
};

export default App;
