import React, { useEffect, useState } from 'react'
import { render } from 'react-dom'
import { BrowserRouter, Route, Link, Switch } from 'react-router-dom'

const App = () => {
  const [devices, setDevices] = useState([])

  const request = async () => {
    const result = await fetch('/api/json')
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
                <td>
                  <Link to={`/${key}`}>{key}</Link>
                </td>
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

const Mac = props => {
  const { mac } = props.match.params

  const [device, setDevice] = useState({})

  const request = async () => {
    const result = await fetch(`/api/${mac}`)
    const data = await result.json()
    setDevice(data)
  }

  useEffect(() => {
    request()
  }, [])

  const onSubmit = async event => {
    event.preventDefault()
    console.log('here')
    const result = await fetch(`/api/${mac}`, {
      method: 'POST',
      body: JSON.stringify(device)
    })
    const data = await result.text()
    console.log(data)
  }

  const onChangeIP = event => {
    const IPParameter = {
      ...device.IPParameter,
      IPAddress: event.target.value
    }
    const next = {
      ...device,
      IPParameter
    }
    setDevice(next)
  }

  // make sure we have a device
  if (!device.Source) {
    return null
  }

  return (
    <div>
      <form onSubmit={onSubmit}>
        <label htmlFor="ip">IP address</label>
        <input
          type="text"
          value={device.IPParameter.IPAddress}
          onChange={onChangeIP}
        />
        <button type="submit">Save</button>
      </form>
    </div>
  )
}

const Index = () => {
  return (
    <Switch>
      <Route exact={true} path="/" component={App} />
      <Route path="/:mac" component={Mac} />
    </Switch>
  )
}

render(
  <BrowserRouter>
    <Index />
  </BrowserRouter>,
  document.getElementById('app')
)
