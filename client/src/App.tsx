import { QueryClient,  QueryClientProvider } from '@tanstack/react-query'
import './App.css'
import AppRoutes from './AppRoutes'
import Notify from './shared/components/Notify'
import { WebSocketProvider } from './shared/contexts/WebSocketContext'

const queryClient = new QueryClient()

const App = () => {

  return (
     <QueryClientProvider client={queryClient}>
       <WebSocketProvider>
         <Notify />
         <AppRoutes />
       </WebSocketProvider>
     </QueryClientProvider>
  )
}

export default App
