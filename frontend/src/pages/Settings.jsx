import React, { useState, useEffect } from 'react'

const API_BASE = '/api'

function Settings() {
  const [settings, setSettings] = useState({
    ai_provider: 'minimax',
    minimax_api_key: '',
    minimax_base_url: 'https://api.minimax.chat',
    n8n_base_url: 'http://localhost:5678',
    comfyui_base_url: 'http://localhost:8188',
    ollama_base_url: 'http://localhost:11434',
    oss_endpoint: 'oss-cn-hangzhou.aliyuncs.com',
    oss_access_key: '',
    oss_secret_key: '',
    oss_bucket: 'genvideo'
  })
  const [saved, setSaved] = useState(false)
  const [loading, setLoading] = useState(false)

  useEffect(() => {
    fetch(`${API_BASE}/settings`)
      .then(res => res.json())
      .then(data => {
        if (data.code === 0 && data.data) {
          setSettings({ ...settings, ...data.data })
        }
      })
      .catch(err => console.error('Failed to load settings:', err))
  }, [])

  const handleSave = async () => {
    setLoading(true)
    try {
      const res = await fetch(`${API_BASE}/settings`, {
        method: 'PUT',
        headers: { 'Content-Type': 'application/json' },
        body: JSON.stringify(settings)
      })
      const data = await res.json()
      if (data.code === 0) {
        setSaved(true)
        setTimeout(() => setSaved(false), 3000)
      }
    } catch (err) {
      console.error('Failed to save settings:', err)
    } finally {
      setLoading(false)
    }
  }

  const handleChange = (key, value) => {
    setSettings(prev => ({ ...prev, [key]: value }))
  }

  return (
    <div className="page">
      <h1>设置</h1>
      <p className="subtitle">配置AI提供者和存储服务</p>

      <div className="card">
        <h3>AI提供商</h3>
        <div className="form-group">
          <label>选择AI服务</label>
          <select
            value={settings.ai_provider}
            onChange={e => handleChange('ai_provider', e.target.value)}
          >
            <option value="minimax">Minimax API</option>
            <option value="n8n">自建 n8n + ComfyUI + Ollama</option>
          </select>
        </div>

        {settings.ai_provider === 'minimax' && (
          <>
            <div className="form-group">
              <label>API Key</label>
              <input
                type="password"
                placeholder="输入 Minimax API Key"
                value={settings.minimax_api_key}
                onChange={e => handleChange('minimax_api_key', e.target.value)}
              />
            </div>
            <div className="form-group">
              <label>API Base URL</label>
              <input
                type="text"
                placeholder="https://api.minimax.chat"
                value={settings.minimax_base_url}
                onChange={e => handleChange('minimax_base_url', e.target.value)}
              />
            </div>
          </>
        )}

        {settings.ai_provider === 'n8n' && (
          <>
            <div className="form-group">
              <label>n8n URL</label>
              <input
                type="text"
                placeholder="http://localhost:5678"
                value={settings.n8n_base_url}
                onChange={e => handleChange('n8n_base_url', e.target.value)}
              />
            </div>
            <div className="form-group">
              <label>ComfyUI URL</label>
              <input
                type="text"
                placeholder="http://localhost:8188"
                value={settings.comfyui_base_url}
                onChange={e => handleChange('comfyui_base_url', e.target.value)}
              />
            </div>
            <div className="form-group">
              <label>Ollama URL</label>
              <input
                type="text"
                placeholder="http://localhost:11434"
                value={settings.ollama_base_url}
                onChange={e => handleChange('ollama_base_url', e.target.value)}
              />
            </div>
          </>
        )}

        <h3>阿里云OSS</h3>
        <div className="form-group">
          <label>Endpoint</label>
          <input
            type="text"
            placeholder="oss-cn-hangzhou.aliyuncs.com"
            value={settings.oss_endpoint}
            onChange={e => handleChange('oss_endpoint', e.target.value)}
          />
        </div>
        <div className="form-group">
          <label>Access Key</label>
          <input
            type="password"
            placeholder="输入 Access Key"
            value={settings.oss_access_key}
            onChange={e => handleChange('oss_access_key', e.target.value)}
          />
        </div>
        <div className="form-group">
          <label>Secret Key</label>
          <input
            type="password"
            placeholder="输入 Secret Key"
            value={settings.oss_secret_key}
            onChange={e => handleChange('oss_secret_key', e.target.value)}
          />
        </div>
        <div className="form-group">
          <label>Bucket</label>
          <input
            type="text"
            placeholder="genvideo"
            value={settings.oss_bucket}
            onChange={e => handleChange('oss_bucket', e.target.value)}
          />
        </div>

        <button className="btn-primary" onClick={handleSave} disabled={loading}>
          {loading ? '保存中...' : '保存设置'}
        </button>

        {saved && <div className="success-msg">设置已保存</div>}
      </div>
    </div>
  )
}

export default Settings
