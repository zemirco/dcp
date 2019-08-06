import React, { useEffect, useState } from 'react'
import { render } from 'react-dom'

const App = () => {
  const [devices, setDevices] = useState([])

  const request = async () => {
    const result = await fetch('/json')
    const data = await result.json()
    console.log(data)
    setDevices(data)
  }

  useEffect(() => {
    request()
  }, [])

  return (
    <div>
      <h1>DCP</h1>
      <table>
        <thead>
          <tr>
            <th></th>
            <th>MAC</th>
            <th>IP address</th>
            <th>Name of station</th>
          </tr>
        </thead>
        <tbody>
          {Object.entries(devices).map(([key, value], i) => {
            console.log(key)
            console.log(value)
            return (
              <tr key={i}>
                <td>{i + 1}</td>
                <td>{key}</td>
                <td>{value.IPParameter.IPAddress}</td>

                <td>{value.NameOfStation.NameOfStation}</td>
              </tr>
            )
          })}
        </tbody>
      </table>
    </div>
  )
}

render(<App />, document.getElementById('app'))
