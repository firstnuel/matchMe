import { type ChangeEvent, useState } from 'react'

export const useField = (
  name: string,
  type: string,
  val: string | number = ''
) => {
  const [value, setValue] = useState(val)

const onChange = (
  e:
    | ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>
    | string
    | { target: { value: [number, number] } }
): void => {
  if (typeof e === 'string') {
    setValue(e)
  } else {
    const targetValue = e.target.value
    if (typeof targetValue === 'string' || typeof targetValue === 'number') {
      setValue(targetValue)
    }
  }
}

  const reset = () => setValue('')

  return {
    name,
    value,
    reset,
    type,
    onChange,
    required: true,
  }
}


