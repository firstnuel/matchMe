import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'
import './App.css'
import AppRoutes from './AppRoutes'


const queryClient = new QueryClient()

const App = () => {

  return (
     <QueryClientProvider client={queryClient}>
       <AppRoutes />
     </QueryClientProvider>
  )
}

export default App
