import React, { useState } from 'react'
import axios from 'axios'

const API_BASE = '/api/v1'

function ScriptToVideo() {
  const [inputType, setInputType] = useState('keywords')
  const [input, setInput] = useState('')
  const [style, setStyle] = useState('dramatic')
  const [duration, setDuration] = useState(60)
  const [sceneCount, setSceneCount] = useState(5)
  const [tasks, setTasks] = useState([])
  const [script, setScript] = useState('')
  const [storyboard, setStoryboard] = useState('')

  const handleGenerateScript = async () => {
    try {
      const response = await axios.post(`${API_BASE}/script/generate`, {
        input,
        input_type: inputType,
        style,
        duration,
      })
      const newTask = response.data.data
      setTasks([...tasks, newTask])

      // Poll for result
      const pollTask = async () => {
        const res = await axios.get(`${API_BASE}/task/${newTask.id}`)
        if (res.data.data.status === 'completed') {
          setScript(res.data.data.output)
        } else if (res.data.data.status !== 'failed') {
          setTimeout(pollTask, 2000)
        }
      }
      pollTask()
    } catch (error) {
      console.error('Script generation failed:', error)
    }
  }

  const handleGenerateStoryboard = async () => {
    if (!script) return
    try {
      const response = await axios.post(`${API_BASE}/script/storyboard`, {
        script,
        scene_count: sceneCount,
      })
      const newTask = response.data.data
      setTasks([...tasks, newTask])

      const pollTask = async () => {
        const res = await axios.get(`${API_BASE}/task/${newTask.id}`)
        if (res.data.data.status === 'completed') {
          setStoryboard(res.data.data.output)
        } else if (res.data.data.status !== 'failed') {
          setTimeout(pollTask, 2000)
        }
      }
      pollTask()
    } catch (error) {
      console.error('Storyboard generation failed:', error)
    }
  }

  const handleGenerateFrames = async () => {
    if (!storyboard) return
    try {
      const response = await axios.post(`${API_BASE}/script/frames`, {
        storyboard,
        style,
      })
      setTasks([...tasks, response.data.data])
    } catch (error) {
      console.error('Frame generation failed:', error)
    }
  }

  const handleGenerateVideo = async () => {
    if (!storyboard) return
    try {
      const response = await axios.post(`${API_BASE}/script/video`, {
        storyboard,
        frames: [],
        duration,
      })
      setTasks([...tasks, response.data.data])
    } catch (error) {
      console.error('Video generation failed:', error)
    }
  }

  return (
    <div className="script-to-video">
      <div className="card">
        <h2>脚本转视频</h2>
        <p style={{ marginBottom: '20px', color: '#666' }}>
          输入关键词、文档或小说，AI 自动生成完整视频
        </p>

        <div className="form-group">
          <label>输入类型</label>
          <select value={inputType} onChange={(e) => setInputType(e.target.value)}>
            <option value="keywords">关键词</option>
            <option value="document">文档</option>
            <option value="novel">小说</option>
          </select>
        </div>

        <div className="form-group">
          <label>输入内容</label>
          <textarea
            value={input}
            onChange={(e) => setInput(e.target.value)}
            placeholder={
              inputType === 'keywords'
                ? '例如：科幻未来城市、机器人觉醒'
                : inputType === 'document'
                ? '粘贴文档内容...'
                : '输入小说内容...'
            }
          />
        </div>

        <div className="form-group">
          <label>视频风格</label>
          <select value={style} onChange={(e) => setStyle(e.target.value)}>
            <option value="dramatic">戏剧化</option>
            <option value="comedy">喜剧</option>
            <option value="documentary">纪录片</option>
            <option value="action">动作</option>
            <option value="romance">浪漫</option>
          </select>
        </div>

        <div className="form-group">
          <label>目标时长（秒）</label>
          <input
            type="number"
            value={duration}
            onChange={(e) => setDuration(parseInt(e.target.value))}
            min={10}
            max={300}
          />
        </div>

        <button className="btn btn-primary" onClick={handleGenerateScript}>
          1. 生成剧本
        </button>
      </div>

      {script && (
        <div className="card">
          <h2>生成的剧本</h2>
          <pre style={{ whiteSpace: 'pre-wrap', marginBottom: '20px' }}>{script}</pre>

          <div className="form-group">
            <label>分镜数量</label>
            <input
              type="number"
              value={sceneCount}
              onChange={(e) => setSceneCount(parseInt(e.target.value))}
              min={1}
              max={20}
            />
          </div>

          <button className="btn btn-primary" onClick={handleGenerateStoryboard}>
            2. 生成分镜
          </button>
        </div>
      )}

      {storyboard && (
        <div className="card">
          <h2>生成的分镜</h2>
          <pre style={{ whiteSpace: 'pre-wrap', marginBottom: '20px' }}>{storyboard}</pre>

          <div style={{ display: 'flex', gap: '10px' }}>
            <button className="btn btn-primary" onClick={handleGenerateFrames}>
              3. 生成首尾帧
            </button>
            <button className="btn btn-secondary" onClick={handleGenerateVideo}>
              4. 直接生成视频
            </button>
          </div>
        </div>
      )}

      {tasks.length > 0 && (
        <div className="card">
          <h2>任务进度</h2>
          <div className="task-list">
            {tasks.map((task, index) => (
              <div key={index} className="task-item">
                <div>
                  <strong>任务 #{task.id}</strong>
                  <p style={{ fontSize: '12px', color: '#666' }}>{task.type}</p>
                </div>
                <div className="progress-bar">
                  <div
                    className="progress-fill"
                    style={{ width: `${task.progress}%` }}
                  />
                </div>
                <span className={`task-status status-${task.status}`}>
                  {task.status}
                </span>
              </div>
            ))}
          </div>
        </div>
      )}
    </div>
  )
}

export default ScriptToVideo
