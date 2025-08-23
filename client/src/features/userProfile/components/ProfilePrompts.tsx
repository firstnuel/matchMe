import React, { useState, type Dispatch, type SetStateAction } from 'react';
import { promptQuestions } from '../types/user';
import {  type Prompt } from "../../../shared/types/user";

interface ProfilePromptsProps {
  initialPrompts?: Prompt[];
  prompts: Prompt[];
  setPrompts: Dispatch<SetStateAction<Prompt[]>>;
}


const ProfilePrompts: React.FC<ProfilePromptsProps> = ({
  initialPrompts = [
    {
      question: "What's your ideal Sunday like?",
      answer: "Sleeping in, making pancakes, and exploring a new hiking trail with good company."
    },
    {
      question: "Two truths and a lie",
      answer: "I've been skydiving, I can speak three languages, I've never broken a bone."
    },
    {
      question: "My biggest fear is...",
      answer: "Running out of coffee on a Monday morning!"
    },
  ],
  prompts,
  setPrompts
}) => {
  // Initialize state with only prompts that have non-empty answers
  const [localPrompts, setLocalPrompts] = useState<Prompt[]>(
    initialPrompts.filter(prompt => prompt.answer.trim() !== "")
  );

  // Use props.prompts if provided and non-empty, otherwise use localPrompts
  const activePrompts = prompts.length > 0 ? prompts : localPrompts;
  const setActivePrompts = setPrompts || setLocalPrompts;

  const handleAnswerChange = (index: number, newAnswer: string) => {
    setActivePrompts(activePrompts.map((prompt, i) =>
      i === index ? { ...prompt, answer: newAnswer } : prompt
    ));
  };

  const removePrompt = (index: number) => {
    setActivePrompts(activePrompts.filter((_, i) => i !== index));
  };

  const addPrompt = () => {
    const availableQuestions = promptQuestions.filter(
      q => !activePrompts.some(p => p.question === q)
    );
    
    const randomQuestion = availableQuestions[Math.floor(Math.random() * availableQuestions.length)] || "New prompt question";

    const newPrompt: Prompt = {
      id: Date.now().toString(), // Optional ID for new prompts
      question: randomQuestion,
      answer: ""
    };
    setActivePrompts([...activePrompts, newPrompt]);
  };

  return (
    <div>
      {activePrompts.map((prompt, index) => (
        <div key={prompt.id ?? index} className="prompt-item">
          <div className="prompt-question">{prompt.question}</div>
          <textarea
            className="form-textarea"
            maxLength={500}
            value={prompt.answer}
            onChange={(e) => handleAnswerChange(index, e.target.value)}
            rows={4}
          />
          <div className="character-count">
            {prompt.answer.length}/500
          </div>
          <div className="prompt-actions">
            <button 
              className="remove-prompt-btn"
              onClick={() => removePrompt(index)}
            >
              Remove
            </button>
          </div>
        </div>
      ))}

      <button className="add-prompt-btn" onClick={addPrompt}>
        + Add Another Prompt
      </button>
    </div>
  );
};

export default ProfilePrompts;