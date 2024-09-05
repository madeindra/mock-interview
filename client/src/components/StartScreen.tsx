import React from 'react';
import { useNavigate } from 'react-router-dom';
import Navbar from './Navbar';
import { useInterviewStore } from '../store';

interface StartScreenProps {
  backendHost: string;
  setError: (error: string | null) => void;
}

const languageOptions = [
  { name: "English", code: "en-US" },
  { name: "Bahasa Indonesia", code: "id-ID" },
];

const StartScreen: React.FC<StartScreenProps> = ({ backendHost, setError }) => {
  const { role, skills, language, messages, setIsIntroDone, setMessages, setRole, setSkills, setLanguage, setInterviewId, setInterviewSecret, setInitialAudio, setInitialText, setInitialSSML } = useInterviewStore();

  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const skillsArray = skills.split(',').map(skill => skill.trim());

    setMessages([]);

    navigate('/processing');

    try {
      const response = await fetch(`${backendHost}/chat/start`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ role, skills: skillsArray, language }),
      });

      const data = await response.json();

      if (response.ok && data.data) {
        setInterviewId(data.data?.id);
        setInterviewSecret(data.data?.secret);
        setInitialAudio(data.data?.audio);
        setInitialText(data.data?.text);
        setInitialSSML(data.data?.ssml);
        setLanguage(data.data?.language);

        setMessages([{ text: data.data?.text, isUser: false, isAnimated: true }]);
        setIsIntroDone(false);

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

  const handleForward = () => {
    navigate('/chat');
  };

  return (
    <div className="flex flex-col h-screen bg-[#1E1E2E] text-white">
      {messages.length > 0 && (
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
                onChange={(e) => setRole(e.target.value)}
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
                onChange={(e) => setSkills(e.target.value)}
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
                onChange={(e) => setLanguage(e.target.value)}
                className="w-full p-3 bg-[#3A3A4E] text-white border border-[#4A4A5E] rounded-lg focus:outline-none focus:ring-2 focus:ring-[#3E64FF]"
                required
              >
                {languageOptions.map((lang) => (
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