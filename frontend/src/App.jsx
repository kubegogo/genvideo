import React from 'react'
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'
import VideoRepurposing from './pages/VideoRepurposing'
import ScriptToVideo from './pages/ScriptToVideo'
import Settings from './pages/Settings'
import './App.css'

function App() {
  return (
    <Router>
      <div className="app">
        <nav className="sidebar">
          <h1>GenVideo</h1>
          <ul>
            <li><Link to="/">视频搬运</Link></li>
            <li><Link to="/script-to-video">脚本转视频</Link></li>
            <li><Link to="/settings">设置</Link></li>
          </ul>
        </nav>
        <main className="content">
          <Routes>
            <Route path="/" element={<VideoRepurposing />} />
            <Route path="/script-to-video" element={<ScriptToVideo />} />
            <Route path="/settings" element={<Settings />} />
          </Routes>
        </main>
      </div>
    </Router>
  )
}

export default App
