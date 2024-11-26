// import {useState} from 'react';
// import logo from './assets/images/logo-universal.png';
// import './App.css';
// import {Greet} from "../wailsjs/go/main/App";

// function App() {
//     const [resultText, setResultText] = useState("Please enter your name below 👇");
//     const [name, setName] = useState('');
//     const updateName = (e) => setName(e.target.value);
//     const updateResultText = (result) => setResultText(result);

//     function greet() {
//         Greet(name).then(updateResultText);
//     }

//     return (
//         <div id="App">
//             <img src={logo} id="logo" alt="logo"/>
//             <div id="result" className="result">{resultText}</div>
//             <div id="input" className="input-box">
//                 <input id="name" className="input" onChange={updateName} autoComplete="off" name="input" type="text"/>
//                 <button className="btn" onClick={greet}>Greet</button>
//             </div>
//         </div>
//     )
// }

import { useState } from "react";

import {SaveFile, SelectFile, ConvertExcelToICal} from "../wailsjs/go/main/App";

function App() {
  const [inputPath, setInputPath] = useState("");
  const [outputPath, setOutputPath] = useState("");
  const [message, setMessage] = useState("");

  const selectFile = async () => {
    const path = await SelectFile();
    setInputPath(path);
  };

  const saveFile = async () => {
    const path = await SaveFile();
    setOutputPath(path);
  };

  const convert = async () => {
    if (!inputPath || !outputPath) {
      setMessage("Please select input and output paths.");
      return;
    }
    try {
      const result = await ConvertExcelToICal(inputPath, outputPath);
      setMessage(result);
    } catch (err) {
      setMessage(`Error: ${err.message}`);
    }
  };

  return (
    <div className="App">
      <h1>Excel to iCal Converter</h1>
      <button onClick={selectFile}>Select Excel File</button>
      <p>{inputPath}</p>
      <button onClick={saveFile}>Save iCal File</button>
      <p>{outputPath}</p>
      <button onClick={convert}>Convert</button>
      <p>{message}</p>
    </div>
  );
}

export default App;


// export default App
