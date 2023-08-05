import { BrowserRouter as Router, Routes, Route } from 'react-router-dom'
import {Login} from './credential-base/Login'

const App = () => {
  return (
    <Router>
      <Routes>
        <Route path="/login">
          <Login />
        </Route>
        {/* 他のルートをここに追加できます */}
      </Routes>
    </Router>
  )
}

export default App
