import React, { useState } from 'react'
import axios from 'axios'

const API = '/api/v1'

export default function VideoGenerate() {
  const [input, setInput] = useState('')
  const [inputType, setInputType] = useState('keywords')
  const [style, setStyle] = useState('documentary')
  const [duration, setDuration] = useState(60)
  const [aspectRatio, setAspectRatio] = useState('16:9')
  const [music, setMusic] = useState('cinematic')
  const [taskId, setTaskId] = useState(null)
  const [task, setTask] = useState(null)

  const handleGenerate = async () => {
    try {
      const res = await axios.post(`${API}/video/generate`, {
        input,
        input_type: inputType,
        style,
        duration,
        aspect_ratio: aspectRatio,
        music,
      })
      const t = res.data.data
      setTaskId(t.id)
      setTask(t)

      // 轮询状态
      const poll = async () => {
        try {
          const r = await axios.get(`${API}/task/${t.id}`)
          setTask(r.data.data)
          if (r.data.data.status !== 'completed' && r.data.data.status !== 'failed') {
            setTimeout(poll, 2000)
          }
        } catch (err) {
          console.error('Poll error:', err)
        }
      }
      poll()
    } catch (err) {
      console.error('Generate error:', err)
    }
  }

  return (
    <div>
      <div className="card">
        <h2>AI视频生成</h2>
        <p style={{ color: '#666', marginBottom: 20 }}>
          输入关键词/文档/小说，AI自动生成剧本→分镜→首尾帧→视频
        </p>

        <div className="form-group">
          <label>输入类型</label>
          <select value={inputType} onChange={e => setInputType(e.target.value)}>
            <option value="keywords">关键词</option>
            <option value="document">文档</option>
            <option value="novel">小说</option>
          </select>
        </div>

        <div className="form-group">
          <label>输入内容</label>
          <textarea
            value={input}
            onChange={e => setInput(e.target.value)}
            placeholder={
              inputType === 'keywords' ? '例如：未来城市、科幻、机器人'
                : inputType === 'document' ? '输入文档内容...'
                : '输入小说内容...'
            }
          />
        </div>

        <div className="form-group">
          <label>视频风格</label>
          <select value={style} onChange={e => setStyle(e.target.value)}>
            <option value="documentary">纪录片</option>
            <option value="dramatic">戏剧化</option>
            <option value="comedy">喜剧</option>
            <option value="action">动作</option>
          </select>
        </div>

        <div className="form-group">
          <label>时长（秒）</label>
          <input type="number" value={duration} onChange={e => setDuration(+e.target.value)} min={10} max={300} />
        </div>

        <div className="form-group">
          <label>画面比例</label>
          <select value={aspectRatio} onChange={e => setAspectRatio(e.target.value)}>
            <option value="16:9">横屏 (16:9)</option>
            <option value="9:16">竖屏 (9:16)</option>
            <option value="1:1">方形 (1:1)</option>
          </select>
        </div>

        <button className="btn btn-primary" onClick={handleGenerate}>
          生成视频
        </button>
      </div>

      {task && (
        <div className="card">
          <h2>任务 #{task.id}</h2>
          <span className={`task-status status-${task.status}`}>{task.status}</span>
          <div className="progress-bar">
            <div className="progress-fill" style={{ width: `${task.progress}%` }} />
          </div>
          {task.status === 'completed' && task.output && (
            <p style={{ marginTop: 15 }}>视频路径: {task.output}</p>
          )}
        </div>
      )}
    </div>
  )
}
