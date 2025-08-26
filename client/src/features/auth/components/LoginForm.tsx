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
import { useUIStore } from '../../../shared/hooks/uiStore'
import { useState, type FormEvent } from 'react'


const LoginForm = () => {
  const [showPassword, setShowPassword] = useState(false)
  const { reset: emailReset, ...email } = useField('email', 'email', '')
  const { reset: passwordReset, ...password } = useField('password', showPassword? 'text': 'password')
  const { errorMsg, infoMsg, setInfo, setError, clearMsgs } = useUIStore()
  const { setAuthToken } = useAuthStore()
  const navigate = useNavigate()
  const mutation = useMutation({
    mutationFn: loginUser,
    onSuccess: (data) => {
        if (data && 'token' in data) {
            setInfo("User Login Successful")
            setAuthToken(data?.token)
            emailReset()
            passwordReset()
            navigate("/")
        }

        if (data && ('error' in data || 'details' in data)) {
            setError(String(data?.details ?? 'An error occurred'))
        }
    }
  })

 const loginFn = (e: FormEvent): void => {
    e.preventDefault()
    clearMsgs()

    const formData: LoginData = {
        email: email.value as string,
        password: password.value as string,
    }

    mutation.mutate(formData)
 }

  const togglePasswordVisibility = () => {
    setShowPassword(!showPassword)
    }

  return ( 
    <div className='container-fluid' >       
        <Container className='form-container'>
            <AppNameTag />
            <Form onSubmit={loginFn} className="d-grid gap-2">
            <div className={errorMsg? 'error': 'info'}>
                {errorMsg || infoMsg || 'Log in to your account'}
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
                    <Form.Control size="lg" {...password} placeholder="Password" autoComplete='new-password'/>
                    <InputGroup.Text onClick={togglePasswordVisibility} style={{ cursor: "pointer" }}>
                    {showPassword? 
                        <Icon icon="mdi:eye-off" className="icon" />
                        : <Icon icon="mdi:eye" className="icon" />
                    }
                    </InputGroup.Text>
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