import Container from 'react-bootstrap/Container'
import Button from 'react-bootstrap/Button'
import Form from 'react-bootstrap/Form'
import { NavLink, useNavigate } from 'react-router-dom'
import { InputGroup } from 'react-bootstrap'
import AppNameTag from '../../../shared/components/AppNameTag'
import { Icon } from "@iconify/react/dist/iconify.js";
import { useField } from '../../../shared/hooks/useField';
import { useMutation } from "@tanstack/react-query";
import { loginUser } from '../api/authApi'
import { type LoginData } from '../types/auth'
import '../styles.css'
import { useAuthStore } from '../hooks/authStore'
import { useState, useEffect, type FormEvent } from 'react'


const LoginForm = () => {
  const { reset: emailReset, ...email } = useField('email', 'email', '')
  const { reset: passwordReset, ...password } = useField('password', 'password')
  const [errorMsg, setErrorMsg] = useState<string>("")
  const { setAuthToken } = useAuthStore()
  const navigate = useNavigate()
  const mutation = useMutation({
    mutationFn: loginUser,
    onSuccess: (data) => {
        if (data && 'token' in data) {
            setAuthToken(data?.token)
            emailReset()
            passwordReset()
            navigate("/")
        }

        if (data && ('error' in data || 'details' in data)) {
            setErrorMsg(String(data?.details ?? 'An error occurred'))
        }
    }
  })

  useEffect(() => {
    if (errorMsg) {
        const timer = setTimeout(() => {
            setErrorMsg("")
        }, 5000)

        return () => clearTimeout(timer)
    }
  }, [errorMsg])


 const loginFn = (e: FormEvent): void => {
    e.preventDefault()
    setErrorMsg("")

    const formData: LoginData = {
        email: email.value as string,
        password: password.value as string,
    }

    mutation.mutate(formData)
 }

  return ( 
    <div className='container-fluid' >       
        <Container className='form-container'>
            <AppNameTag />
            <Form onSubmit={loginFn} className="d-grid gap-2">
            <div className={errorMsg? 'error': 'info'}>
                {errorMsg? errorMsg : 'Log in to your account'}
            </div>
                <InputGroup className="mb-3">
                    <InputGroup.Text> 
                        <Icon icon="mdi:user" className="icon" />
                    </InputGroup.Text>
                    <Form.Control size="lg" { ...email } placeholder="Email" autoComplete='current-email'/>
                </InputGroup>

                <InputGroup className="mb-3">
                    <InputGroup.Text>
                    <Icon icon="mdi:lock" className="icon" />
                    </InputGroup.Text>
                    <Form.Control size="lg" {...password} placeholder="Password"  autoComplete='current-password'/>
                </InputGroup>

                <Button variant="primary" type="submit"  size="lg" className='formButton'>
                    {mutation.isPending? 'Loading...' : 'Login'}
                </Button >
            </Form>
            <div className='register' style={{ marginTop: "1em" }}>
                <p>No account yet? <NavLink to='/register' className='register-link'>Register</NavLink></p>
        </div>
      </Container>
    </div>

    )
}

export default LoginForm;