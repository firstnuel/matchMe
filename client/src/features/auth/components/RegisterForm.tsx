import Container from 'react-bootstrap/Container'
import Button from 'react-bootstrap/Button'
import Form from 'react-bootstrap/Form'
import { NavLink } from 'react-router-dom'
import { InputGroup } from 'react-bootstrap'
import AppNameTag from '../../../shared/components/AppNameTag'
import { Icon } from "@iconify/react/dist/iconify.js";
import { useField } from '../../../shared/hooks/useField';
import { registerUser } from '../api/authApi'
import { useMutation } from "@tanstack/react-query";
import { useAuthStore } from '../hooks/authStore'
import { useUIStore } from '../../../shared/hooks/uiStore'
import { calculateAge } from '../../../shared/utils/ageHelper'
import { useNavigate } from 'react-router-dom'
import { useState, useEffect, type FormEvent } from 'react'
import '../styles.css'
import type { RegisterData } from '../types/auth'


const RegisterForm = () => {
  const { reset: emailReset, ...email } = useField('email', 'email', '')
  const { reset: passwordReset, ...password } = useField('password', 'password')
  const { reset: firstNameReset, ...firstName } = useField('text', 'text', '')
  const { reset: lastNameReset, ...lastName } = useField('text', 'text', '')
  const { reset: genderReset, ...gender } = useField('text', 'text', '')
  const { reset: ageReset, ...age } = useField('date', 'date', '')
  const navigate = useNavigate()
  const { setView } = useUIStore()
  const [errorMsg, setErrorMsg] = useState<string>("")
  const [infoMsg, setInfoMsg] = useState<string>("")

    const { setAuthToken } = useAuthStore()
  const mutation = useMutation({
    mutationFn: registerUser,
    onSuccess: (data) => {
        if (data && 'token' in data) {
            setInfoMsg("User Registration Successful")
            setAuthToken(data?.token)
            emailReset()
            passwordReset()
            firstNameReset()
            lastNameReset()
            genderReset()
            ageReset()
            navigate("/")
            setView("home")
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
     
  const registerFn = (e: FormEvent): void => {
      e.preventDefault()
      setErrorMsg("")

      const formData: RegisterData = {
          email: email.value as string,
          password: password.value as string,
          first_name: firstName.value as string,
          last_name: lastName.value as string,
          gender: gender.value as "male" | "female" | "non_binary" | "prefer_not_to_say",
          age: calculateAge(age.value),

      }
      mutation.mutate(formData)
  }
  
  return ( 
    <div className='container-fluid' >       
      <Container className='form-container'>
        <AppNameTag />
        <Form onSubmit={registerFn} className="d-grid gap-2">
          <div className={errorMsg? 'error': 'info'}>
            {errorMsg || infoMsg || 'Create your account'}
          </div>
          
          <InputGroup className="mb-3">
            <InputGroup.Text> 
              <Icon icon="mdi:account" className="icon" />
            </InputGroup.Text>
            <Form.Control size="lg" { ...firstName } placeholder="First Name" autoComplete='given-name'/>
          </InputGroup>

          <InputGroup className="mb-3">
            <InputGroup.Text> 
              <Icon icon="mdi:account" className="icon" />
            </InputGroup.Text>
            <Form.Control size="lg" { ...lastName } placeholder="Last Name" autoComplete='family-name'/>
          </InputGroup>

          <InputGroup className="mb-3">
            <InputGroup.Text> 
              <Icon icon="mdi:email" className="icon" />
            </InputGroup.Text>
            <Form.Control size="lg" { ...email } placeholder="Email" autoComplete='email'/>
          </InputGroup>

          <InputGroup className="mb-3">
            <InputGroup.Text>
              <Icon icon="mdi:lock" className="icon" />
            </InputGroup.Text>
            <Form.Control size="lg" {...password} placeholder="Password" autoComplete='new-password'/>
          </InputGroup>

          <InputGroup className="mb-3">
            <InputGroup.Text>
              <Icon icon="mdi:gender-male-female" className="icon" />
            </InputGroup.Text>
            <Form.Select size="lg" {...gender}>
              <option value="">Select Gender</option>
              <option value="male">Male</option>
              <option value="female">Female</option>
              <option value="non_binary">Non Binary</option>
              <option value="prefer_not_to_say">Undisclosed</option>
            </Form.Select>
          </InputGroup>

        <InputGroup className="mb-3">
          <InputGroup.Text>
            <Icon icon="mdi:calendar" className="icon" />
          </InputGroup.Text>
          <Form.Control 
            size="lg" 
            {...age} 
            type="date" 
            placeholder="Date of Birth" 
            />
        </InputGroup>
        <div>
              <input style={{ margin: "0.4rem" }} type="checkbox" id="age-check" required/>
              <label htmlFor="age-check">
                By clicking this I accept that I am 18 and older
              </label>
          </div>
          <Button disabled={mutation.isPending}
            variant="primary" type="submit" size="lg" className='formButton'>
            Register
          </Button >
        </Form>
        <div className='register' style={{ marginTop: "1em" }}>
          <p>Already have an account? <NavLink to='/login' className='register-link'>Login</NavLink></p>
        </div>
      </Container>
    </div>
  )
}

export default RegisterForm;