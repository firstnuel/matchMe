import { type ChangeEvent, useState } from "react"

export const useSelect = (initialValue: string = "") => {
  const [value, setValue] = useState(initialValue)

  const onChange = (
  e:
    | ChangeEvent<HTMLInputElement | HTMLSelectElement | HTMLTextAreaElement>
    | string
    | { target: { value: [number, number] } }
): void => {
    if (typeof e === "string") {
      setValue(e) // allows programmatic updates
    } else {
      setValue(String(e.target.value)) // updates from <select> UI
    }
  }

  const reset = () => setValue(initialValue)

  return { value, onChange, reset }
}
