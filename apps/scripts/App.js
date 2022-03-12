import logo from "./logo.svg";
import "./App.css";
import io from "socket.io-client";
import { useEffect } from "react";
export const socket = io("http://127.0.0.1:65432");

function App() {
  useEffect(() => {
    console.log("mounted");
    socket.on("testResponse", (resp) => {
      console.log(resp);
    });
  }, []);
  const handleTest = () => {
    console.log("test");
    socket.emit("test", { data: "teeeeeeest" });
  };

  const handleSetPosition = () => {
    console.log("position");
    socket.emit("position", { data: "teeeeeeest" });
  };
  return (
    <div className="App">
      <header className="App-header">
        <img src={logo} className="App-logo" alt="logo" />
        <button onClick={handleTest}>Test connection</button>
        <button onClick={handleSetPosition}>Set Position</button>
      </header>
    </div>
  );
}

export default App;
