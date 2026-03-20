import React from 'react'
import { BrowserRouter as Router, Routes, Route, Link } from 'react-router-dom'
import VideoRepurposing from './pages/VideoRepurposing'
import ScriptToVideo from './pages/ScriptToVideo'
import Settings from './pages/Settings'

function App() {
  return (
    <Router>
      <div className="app">
        <nav className="navbar">
          <div className="nav-brand">GenVideo</div>
          <div className="nav-links">
            <Link to="/">视频搬运</Link>
            <Link to="/script-to-video">脚本转视频</Link>
            <Link to="/settings">设置</Link>
          </div>
        </nav>
        <main className="main-content">
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
