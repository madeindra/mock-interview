import React from 'react';
import { BrowserRouter as Router, Routes, Route } from 'react-router-dom';
import StartScreen from './components/StartScreen';
import ProcessingScreen from './components/ProcessingScreen';
import ChatScreen from './components/ChatScreen';

const App: React.FC = () => {
  return (
    <Router>
      <div className="min-h-screen bg-gray-100">
        <Routes>
          <Route path="/" element={<StartScreen />} />
          <Route path="/processing" element={<ProcessingScreen />} />
          <Route path="/chat" element={<ChatScreen />} />
        </Routes>
      </div>
    </Router>
  );
};

export default App;