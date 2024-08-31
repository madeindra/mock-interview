import React, { useState } from 'react';
import { useNavigate } from 'react-router-dom';

const StartScreen: React.FC = () => {
  const [role, setRole] = useState('');
  const [skills, setSkills] = useState('');
  const navigate = useNavigate();

  const handleSubmit = async (e: React.FormEvent) => {
    e.preventDefault();
    const skillsArray = skills.split(',').map(skill => skill.trim());
    
    try {
      const response = await fetch(`${import.meta.env.VITE_BACKEND_URL}/chat/start`, {
        method: 'POST',
        headers: {
          'Content-Type': 'application/json',
        },
        body: JSON.stringify({ role, skills: skillsArray }),
      });

      if (response.ok) {
        const data = await response.json();
        // Store id and secret securely (e.g., in localStorage or state management solution)
        localStorage.setItem('interviewId', data.id);
        localStorage.setItem('interviewSecret', data.secret);
        
        // Store audio and text for the chat screen
        localStorage.setItem('initialAudio', data.audio);
        localStorage.setItem('initialText', data.text);
        
        navigate('/chat');
      } else {
        console.error('Failed to start interview');
      }
    } catch (error) {
      console.error('Error starting interview:', error);
    }
  };

  return (
    <div className="container mx-auto mt-10 p-4">
      <form onSubmit={handleSubmit} className="max-w-md mx-auto">
        <div className="mb-4">
          <label htmlFor="role" className="block mb-2">Role</label>
          <input
            type="text"
            id="role"
            value={role}
            onChange={(e) => setRole(e.target.value)}
            placeholder="e.g. Backend Engineer"
            className="w-full p-2 border rounded"
            required
          />
        </div>
        <div className="mb-4">
          <label htmlFor="skills" className="block mb-2">Skills</label>
          <textarea
            id="skills"
            value={skills}
            onChange={(e) => setSkills(e.target.value)}
            placeholder="e.g. HTTP, Golang, REST API, AWS"
            className="w-full p-2 border rounded"
            required
          ></textarea>
        </div>
        <button type="submit" className="w-full bg-blue-500 text-white p-2 rounded">
          Start Interview
        </button>
      </form>
    </div>
  );
};

export default StartScreen;