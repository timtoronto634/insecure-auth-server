import React, { useState } from 'react'

export function Login() {
  const [username, setUsername] = useState<string>('')
  const [password, setPassword] = useState<string>('')

  const handleSubmit = (event: React.FormEvent) => {
    event.preventDefault()
    // ここでユーザー名とパスワードを使ってログイン処理を行います
    console.log(
      `Logging in with username: ${username} and password: ${password}`,
    )
  }

  return (
    <div>
      <h2>Login</h2>
      <form onSubmit={handleSubmit}>
        <label>
          Username:
          <input
            type="text"
            value={username}
            onChange={(e) => setUsername(e.target.value)}
          />
        </label>
        <br />
        <label>
          Password:
          <input
            type="password"
            value={password}
            onChange={(e) => setPassword(e.target.value)}
          />
        </label>
        <br />
        <input type="submit" value="Login" />
      </form>
    </div>
  )
}
