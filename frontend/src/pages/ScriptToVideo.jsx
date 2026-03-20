import React, { useState } from 'react'

const API_BASE = '/api'

function ScriptToVideo() {
  const [inputType, setInputType] = useState('keywords')
  const [input, setInput] = useState('')
  const [style, setStyle] = useState('auto')
  const [duration, setDuration] = useState(60)
  const [aspectRatio, setAspectRatio] = useState('16:9')
  const [loading, setLoading] = useState(false)
  const [taskId, setTaskId] = useState(null)
  const [progress, setProgress] = useState(0)
  const [result, setResult] = useState(null)
  const [error, setError] = useState(null)

  const handleGenerate = async () => {
    if (!input) {
      setError('请输入内容')
      return
    }
    setLoading(true)
    setError(null)
    setProgress(0)
    try {
      const res = await fetch(`${API_BASE}/video/generate`, {
        method: 'POST',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify({
          input,
          input_type: inputType,
          style,
          duration,
          aspect_ratio: aspectRatio
        })
      })
      const data = await res.json()
      if (data.code === 0) {
        setTaskId(data.data.task_id)
        pollTaskStatus(data.data.task_id)
      } else {
        setError(data.msg)
      }
    } catch (err) {
      setError('请求失败: ' + err.message)
    } finally {
      setLoading(false)
    }
  }

  const pollTaskStatus = async (id) => {
    const interval = setInterval(async () => {
      try {
        const res = await fetch(`${API_BASE}/task/${id}`)
        const data = await res.json()
        if (data.code === 0) {
          const task = data.data
          setProgress(task.progress || 0)
          if (task.status === 'completed') {
            setResult(task.output)
            clearInterval(interval)
          } else if (task.status === 'failed') {
            setError(task.error)
            clearInterval(interval)
          }
        }
      } catch (err) {
        console.error('Poll error:', err)
      }
    }, 3000)
  }

  return (
    <div className="page">
      <h1>脚本转视频</h1>
      <p className="subtitle">输入关键词/文档/小说 → AI生成剧本 → 自动剪辑配乐 → 生成原创视频</p>

      <div className="card">
        <div className="form-group">
          <label>输入类型</label>
          <select value={inputType} onChange={e => setInputType(e.target.value)}>
            <option value="keywords">关键词</option>
            <option value="text">文本/文章</option>
            <option value="document">文档</option>
            <option value="novel">小说</option>
          </select>
        </div>

        <div className="form-group">
          <label>输入内容</label>
          <textarea
            rows={6}
            placeholder={inputType === 'keywords' ? '例如：日出、云海、山峰' : '输入文本内容...'}
            value={input}
            onChange={e => setInput(e.target.value)}
          />
        </div>

        <div className="form-row">
          <div className="form-group">
            <label>视频风格</label>
            <select value={style} onChange={e => setStyle(e.target.value)}>
              <option value="auto">自动匹配</option>
              <option value="documentary">纪录片</option>
              <option value="vlog">Vlog</option>
              <option value="cinematic">电影感</option>
            </select>
          </div>

          <div className="form-group">
            <label>时长 (秒)</label>
            <input
              type="number"
              min={15}
              max={180}
              value={duration}
              onChange={e => setDuration(parseInt(e.target.value))}
            />
          </div>

          <div className="form-group">
            <label>比例</label>
            <select value={aspectRatio} onChange={e => setAspectRatio(e.target.value)}>
              <option value="16:9">16:9 横版</option>
              <option value="9:16">9:16 竖版</option>
              <option value="1:1">1:1 方形</option>
            </select>
          </div>
        </div>

        {loading && (
          <div className="progress-bar">
            <div className="progress-fill" style={{ width: `${progress}%` }} />
          </div>
        )}

        <button className="btn-primary" onClick={handleGenerate} disabled={loading}>
          {loading ? '生成中...' : '开始生成视频'}
        </button>

        {error && <div className="error-msg">{error}</div>}

        {result && (
          <div className="result">
            <h3>视频生成完成!</h3>
            <a href={result} target="_blank" rel="noopener noreferrer">下载视频</a>
          </div>
        )}
      </div>
    </div>
  )
}

export default ScriptToVideo
