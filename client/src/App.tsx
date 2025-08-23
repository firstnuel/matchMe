import {
  QueryClient,
  QueryClientProvider,
} from '@tanstack/react-query'
import './App.css'
import AppRoutes from './AppRoutes'
import Notify from './shared/components/Notify'

const queryClient = new QueryClient()

const App = () => {

  return (
     <QueryClientProvider client={queryClient}>
      <Notify />
       <AppRoutes />
     </QueryClientProvider>
  )
}

export default App
