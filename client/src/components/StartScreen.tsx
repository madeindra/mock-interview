import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Navbar from './Navbar';

interface StartScreenProps {
  setError: (error: string | null) => void;
}

const StartScreen: React.FC<StartScreenProps> = ({ setError }) => {
  const tempRole = sessionStorage.getItem('role');
  const tempSkills = sessionStorage.getItem('skills');

  const [role, setRole] = useState(tempRole || '');
  const [skills, setSkills] = useState(tempSkills || '');
  const [hasMessages, setHasMessages] = useState(false);

  const navigate = useNavigate();

  useEffect(() => {
    const storedMessages = sessionStorage.getItem('messages');
    setHasMessages(!!storedMessages);
  }, []);

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const skillsArray = skills.split(',').map(skill => skill.trim());

    sessionStorage.removeItem('messages');

    navigate('/processing');

    try {
      const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/chat/start`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ role, skills: skillsArray }),
      });

      const data = await response.json();

      if (response.ok && data.data) {
        sessionStorage.setItem('interviewId', data.data.id);
        sessionStorage.setItem('interviewSecret', data.data.secret);
        sessionStorage.setItem('initialAudio', data.data.audio);
        sessionStorage.setItem('initialText', data.data.text);
        navigate('/chat');
      } else {
        const errorMessage = data.message || 'Failed processing your request, please try again';
        setError(errorMessage);
        navigate('/');
      }
    } catch (error) {
      console.error('Error starting interview:', error);
      setError('Failed processing your request, please try again');
      navigate('/');
    }
  };

  const updateRole = (role: string) => {
    setRole(role);
    sessionStorage.setItem('role', role);
  }

  const updateSkills = (skills: string) => {
    setSkills(skills);
    sessionStorage.setItem('skills', skills);
  }

  const handleForward = () => {
    navigate('/chat');
  };

  return (
    <div className="flex flex-col h-screen">
      {hasMessages && (
        <Navbar
          showBackIcon
          showForwardIcon
          onForward={handleForward}
          disableBack={true}
        />
      )}
      <div className="container mx-auto mt-10 p-4 flex-grow">
        <div className="max-w-md mx-auto bg-dark-surface p-8 rounded-xl shadow-lg">
          <h1 className="text-3xl font-bold mb-6 text-center text-dark-on-surface">Mock Interview</h1>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="role" className="block mb-2 text-dark-on-surface font-semibold">Role</label>
              <input
                type="text"
                id="role"
                value={role}
                onChange={(e) => updateRole(e.target.value)}
                placeholder="e.g. Fullstack Developer"
                className="input-field w-full"
                required
              />
            </div>
            <div>
              <label htmlFor="skills" className="block mb-2 text-dark-on-surface font-semibold">Skills</label>
              <textarea
                id="skills"
                value={skills}
                onChange={(e) => updateSkills(e.target.value)}
                placeholder="e.g. Javascript, Typescript, REST API"
                className="input-field w-full h-32"
                required
              ></textarea>
            </div>
            <button type="submit" className="btn-primary w-full text-dark-on-surface">
              Start Interview
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default StartScreen;