import React, { useState, useEffect } from 'react';
import { Message } from '../store';

const useTypingEffect = (message: Message, speed: number = 25) => {
    const [displayedText, setDisplayedText] = useState('');

    useEffect(() => {
        if (message.isAnimated && !message.hasAnimated) {
            let i = 0;
            const timer = setInterval(() => {
                setDisplayedText(message.text.slice(0, i));
                i++;
                if (i > message.text.length) {
                    clearInterval(timer);

                    message.hasAnimated = true;
                }
            }, speed);

            return () => clearInterval(timer);
        } else {
            setDisplayedText(message.text);
        }
    }, [message, speed]);

    return displayedText;
};

interface AnimatedTextProps {
    message: Message
  }

const AnimatedText: React.FC<AnimatedTextProps> = ({ message }) => {
    const displayedText = useTypingEffect(message);
    return <>{displayedText}</>;
};

export default AnimatedText;