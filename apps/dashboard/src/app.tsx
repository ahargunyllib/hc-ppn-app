import { useState } from "react";

export default function App() {
  const [count, setCount] = useState(0);

  return (
    <div className="isolate">
      <h1>Vite + React</h1>
      <div>
        <button onClick={() => setCount((prev) => prev + 1)} type="button">
          count is {count}
        </button>
        <p>
          Edit <code>src/App.tsx</code> and save to test HMR
        </p>
      </div>
      <p>Click on the Vite and React logos to learn more</p>
    </div>
  );
}
