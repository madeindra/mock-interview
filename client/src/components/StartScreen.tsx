import React, { useState, useEffect } from 'react';
import { useNavigate } from 'react-router-dom';
import Navbar from './Navbar';

interface StartScreenProps {
  backendHost: string;
  setError: (error: string | null) => void;
}

const languageOptions = [
  { name: "English", code: "en-US" },
  { name: "Bahasa Indonesia", code: "id-ID" },
];

const StartScreen: React.FC<StartScreenProps> = ({ backendHost, setError }) => {
  const tempRole = sessionStorage.getItem('role');
  const tempSkills = sessionStorage.getItem('skills');
  const tempLanguage = sessionStorage.getItem('language');

  const [role, setRole] = useState(tempRole || '');
  const [skills, setSkills] = useState(tempSkills || '');
  const [language, setLanguage] = useState(tempLanguage || 'en');
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
      const response = await fetch(`${backendHost}/chat/start`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ role, skills: skillsArray, language  }),
      });

      const data = await response.json();

      if (response.ok && data.data) {
        sessionStorage.setItem('interviewId', data.data?.id);
        sessionStorage.setItem('interviewSecret', data.data?.secret);
        sessionStorage.setItem('initialAudio', data.data?.audio);
        sessionStorage.setItem('initialText', data.data?.text);
        sessionStorage.setItem('chatLanguage', data.data?.language);
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

  const updateLanguage = (language: string) => {
    setLanguage(language);
    sessionStorage.setItem('language', language);
  }

  const handleForward = () => {
    navigate('/chat');
  };

  return (
    <div className="flex flex-col h-screen bg-[#1E1E2E] text-white">
      {hasMessages && (
        <Navbar
          backendHost={backendHost}
          showBackIcon
          showForwardIcon
          onForward={handleForward}
          disableBack={true}
        />
      )}
      <div className="container mx-auto mt-10 p-4 flex-grow">
        <div className="max-w-md mx-auto bg-[#2B2B3B] p-8 rounded-xl shadow-lg">
          <h1 className="text-3xl font-bold mb-6 text-center text-white">Mock Interview</h1>
          <form onSubmit={handleSubmit} className="space-y-6">
            <div>
              <label htmlFor="role" className="block mb-2 text-white font-semibold">Role</label>
              <input
                type="text"
                id="role"
                value={role}
                onChange={(e) => updateRole(e.target.value)}
                placeholder="e.g. Fullstack Developer"
                className="w-full p-3 bg-[#3A3A4E] text-white border border-[#4A4A5E] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3E64FF]"
                required
              />
            </div>
            <div>
              <label htmlFor="skills" className="block mb-2 text-white font-semibold">Skills</label>
              <textarea
                id="skills"
                value={skills}
                onChange={(e) => updateSkills(e.target.value)}
                placeholder="e.g. Javascript, Typescript, REST API"
                className="w-full h-32 p-3 bg-[#3A3A4E] text-white border border-[#4A4A5E] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3E64FF]"
                required
              ></textarea>
            </div>
            <div>
              <label htmlFor="language" className="block mb-2 text-white font-semibold">Language</label>
              <select
                id="language"
                value={language}
                onChange={(e) => updateLanguage(e.target.value)}
                className="w-full p-3 bg-[#3A3A4E] text-white border border-[#4A4A5E] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3E64FF]"
                required
              >
                {languageOptions.map((lang) =>  (
                  <option key={lang.code} value={lang.code}>{lang.name}</option>
                ))}
              </select>
            </div>
            <button type="submit" className="w-full p-4 bg-[#3E64FF] text-white font-bold rounded-xl hover:bg-opacity-90 transition-all duration-300">
              Start Interview
            </button>
          </form>
        </div>
      </div>
    </div>
  );
};

export default StartScreen;