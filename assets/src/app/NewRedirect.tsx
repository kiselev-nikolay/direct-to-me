import React from 'react';

import { QuestionIcon } from '@chakra-ui/icons';
import {
  Box,
  Button,
  ButtonGroup,
  Code,
  Divider,
  FormControl,
  FormLabel,
  Input,
  ListItem,
  OrderedList,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverContent,
  PopoverTrigger,
  Switch,
  Tab,
  TabList,
  TabPanel,
  TabPanels,
  Tabs,
  Textarea,
} from '@chakra-ui/react';

import { Redirect } from './API/Redirect';
import RedirectUI from './ui/Redirect';

let ifSaved = (key: string, defaultValue: any) => {
  let rawSaved = localStorage.getItem('NewRedirectFrom');
  if (rawSaved !== null && rawSaved !== undefined) {
    let saved: any = JSON.parse(rawSaved);
    if (saved !== null && saved !== undefined) {
      return saved[key] || defaultValue;
    }
  }
  return defaultValue;
};

interface LabelWithHelpProps {
  title: string;
  placeholder?: string;
  children: React.ReactNode;
}
function LabelWithHelp(props: LabelWithHelpProps) {
  return (
    <Popover trigger="hover">
      <FormLabel>
        {props.title}
        <PopoverTrigger>
          <QuestionIcon mx="10px" mt="-5px" color="gray.300" />
        </PopoverTrigger>
      </FormLabel>
      <PopoverContent mx="10px">
        <PopoverArrow />
        <PopoverBody>{props.children}</PopoverBody>
      </PopoverContent>
    </Popover>
  );
}

interface TextControlProps {
  title: string;
  doc: string | React.ReactNode;
  name: string;
  value?: string;
  placeholder?: string;
  on: (name: string, e: any) => void;
}
function LineControl(props: TextControlProps) {
  return (<>
    <FormControl mb="1rem">
      <LabelWithHelp title={props.title}>{props.doc}</LabelWithHelp>
      <Input name={props.name} placeholder={props.placeholder} value={props.value} onChange={e => props.on(props.name, e)} />
    </FormControl>
  </>);
}
function AreaControl(props: TextControlProps) {
  return (<>
    <FormControl mb="1rem">
      <LabelWithHelp title={props.title}>{props.doc}</LabelWithHelp>
      <Textarea name={props.name} placeholder={props.placeholder} value={props.value} onChange={e => props.on(props.name, e)} resize="vertical" />
    </FormControl>
  </>);
}
interface BoolControlProps extends TextControlProps {
  checked?: boolean;
}
function SwitchControl(props: BoolControlProps) {
  return (<>
    <FormControl display="flex" alignItems="center" mb="1rem">
      <LabelWithHelp title={props.title}>{props.doc}</LabelWithHelp>
      <Switch name={props.name} defaultChecked={props.checked} onChange={e => props.on(props.name, e)} />
    </FormControl>
  </>);
}

interface NewRedirectFromProps { }
interface NewRedirectFromState {
  FromURI: string;
  ToURL: string;
  RedirectAfter: string;
  URLTemplate: string;
  MethodTemplate: string;
  HeadersTemplate: string;
  BodyTemplate: string;

  isAdvanced: boolean;
  step: number;
}
class NewRedirectFrom extends React.Component<NewRedirectFromProps, NewRedirectFromState> {
  constructor(props: NewRedirectFromProps) {
    super(props);
    this.state = {
      FromURI: ifSaved('FromURI', ''),
      ToURL: ifSaved('ToURL', ''),
      RedirectAfter: ifSaved('RedirectAfter', ''),
      URLTemplate: ifSaved('URLTemplate', ''),
      MethodTemplate: ifSaved('MethodTemplate', ''),
      HeadersTemplate: ifSaved('HeadersTemplate', ''),
      BodyTemplate: ifSaved('BodyTemplate', ''),
      isAdvanced: ifSaved('isAdvanced', false),
      step: 0,
    };
  }
  onChange(name: string, e: any) {
    let newState: any = {};
    newState[name] = e.target.value;
    this.setState(newState);
    localStorage.setItem("NewRedirectFrom", JSON.stringify(this.state));
    e.preventDefault();
  }
  onSwitch(name: string, e: any) {
    let newState: any = {};
    newState[name] = e.target.checked;
    this.setState(newState);
    localStorage.setItem("NewRedirectFrom", JSON.stringify(this.state));
    e.preventDefault();
  }
  onSubmit(e: any) {
    let r = new Redirect();
    r.FromURI = this.state.FromURI;
    r.RedirectAfter = this.state.RedirectAfter;
    if (this.state.isAdvanced) {
      r.URLTemplate = this.state.URLTemplate;
      r.MethodTemplate = this.state.MethodTemplate;
      r.HeadersTemplate = this.state.HeadersTemplate;
      r.BodyTemplate = this.state.BodyTemplate;
    } else {
      r.ToURL = this.state.ToURL;
    }
    let p = r.createNew();
    p.then((err) => {
      if (err === null) {
        localStorage.removeItem("NewRedirectFrom");
        localStorage.setItem("page", "index");
        location.reload();
      } else {
        // FIXIT show error modal or snick-bar
        alert(err);
      }
    });
    e.preventDefault();
  }
  nextStep() {
    this.setState({ step: this.state.step + 1 });
  }
  backStep() {
    this.setState({ step: this.state.step - 1 });
  }
  render() {
    return (
      <Tabs isFitted index={this.state.step}>
        <TabList>
          <Tab isDisabled={this.state.step < 0}>Step 1: Redirect</Tab>
          <Tab isDisabled={this.state.step < 1}>Step 2: Background send</Tab>
          <Tab isDisabled={this.state.step < 2}>Step 3: Save</Tab>
        </TabList>
        <TabPanels mt="2rem">
          <TabPanel>
            <Box w="100%">
              <LineControl name="FromURI" value={this.state.FromURI}
                title="From URI" placeholder="to-my-site" doc={<>
                  URI used for that redirect. As example redirect with from URI <Code>to-my-site</Code> can be accessed from
                  <Code>{location.origin + "/to-my-site"}</Code>
                </>}
                on={(name, e) => this.onChange(name, e)} />
              <LineControl name="RedirectAfter" value={this.state.RedirectAfter}
                title="Redirect after" placeholder="https://direct-to-me.com/buy?from=to-my-site" doc="URL to redirect after background data send"
                on={(name, e) => this.onChange(name, e)} />
              <ButtonGroup spacing="3">
                <Button colorScheme="green" onClick={e => this.nextStep()}>Next Step</Button>
              </ButtonGroup>
            </Box>
          </TabPanel>
          <TabPanel>
            <Box w="100%">
              <SwitchControl name="isAdvanced" checked={this.state.isAdvanced}
                title="Advanced mode" doc={<>
                  Toggle advanced mode.
                  <OrderedList>
                    <ListItem>Simple mode: Just send data in body as JSON object to webhook</ListItem>
                    <ListItem>Advanced mode: Template entire HTTP request, URL, Method, Headers, Body</ListItem>
                  </OrderedList>
                </>}
                on={(name, e) => this.onSwitch(name, e)} />
              <Divider mb="2rem" />
              {this.state.isAdvanced ? <>
                <LineControl name="URLTemplate" value={this.state.URLTemplate}
                  title="URL template" placeholder="https://hooks.slack.com/services/T{{secrets.TID}}/B{{secrets.BID}}" doc="Template for URL"
                  on={(name, e) => this.onChange(name, e)} />
                <LineControl name="MethodTemplate" value={this.state.MethodTemplate}
                  title="Method template" placeholder="POST" doc="Template for method. You can set it as static 'GET', 'POST' or any"
                  on={(name, e) => this.onChange(name, e)} />
                <AreaControl name="HeadersTemplate" value={this.state.HeadersTemplate}
                  title="Headers template" placeholder={"Authorization: {{secrets.Token}}\nX-Trace-Id: {{data.trace_id}}\nContent-Type: text/xml"} doc="Template for request headers"
                  on={(name, e) => this.onChange(name, e)} />
                <AreaControl name="BodyTemplate" value={this.state.BodyTemplate}
                  title="Body template" placeholder={'<?xml version="1.0" encoding="UTF-8"?>\n<root>\n\t<text>New lead! Email: {{data.lead_email}}</text>\n</root>'} doc="Template for request body"
                  on={(name, e) => this.onChange(name, e)} />
              </> : <>
                <LineControl name="ToURL" value={this.state.ToURL}
                  title="Send JSON to URL in background" placeholder="https://hooks.slack.com/services/T12345678/B12345678" doc="Specify full URL where to send the data. It can be any API endpoint, webhook, or postback"
                  on={(name, e) => this.onChange(name, e)} />
              </>}
              <ButtonGroup spacing="3">
                <Button onClick={e => this.backStep()}>Previous Step</Button>
                <Button colorScheme="green" onClick={e => this.nextStep()}>Next Step</Button>
              </ButtonGroup>
            </Box>
          </TabPanel>
          <TabPanel>
            <Box w="100%">
              <RedirectUI redirect={this.state} />
              <Divider my="1rem" />
              <ButtonGroup spacing="3">
                <Button onClick={e => this.backStep()}>Previous Step</Button>
                <Button colorScheme="green" onClick={e => this.onSubmit(e)}>Submit</Button>
              </ButtonGroup>
            </Box>
          </TabPanel>
        </TabPanels>
      </Tabs>
    );
  }
}

export default NewRedirectFrom;