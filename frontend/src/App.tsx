import { BrowserRouter, Routes, Route } from 'react-router-dom'
import ProductPage from './components/ProductPage'

function App() {
  return (
    <BrowserRouter>
      <Routes>
        <Route path="/product/:slug" element={<ProductPage />} />
      </Routes>
    </BrowserRouter>
  )
}

export default App
