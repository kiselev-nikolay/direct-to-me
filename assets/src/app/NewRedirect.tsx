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
  Heading,
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
  Tag,
  Textarea,
} from '@chakra-ui/react';

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

interface DisplayTemplateProps {
  v: string;
};
function DisplayTemplate(props: DisplayTemplateProps) {
  let text: Array<React.ReactNode> = [];
  const regex = /{{(\s+)?(.+?)(\s+)?}}/gm;
  let m;
  let last = 0;
  while ((m = regex.exec(props.v)) !== null) {
    if (m.index === regex.lastIndex) {
      regex.lastIndex++;
    }
    text.push(<span key={text.length}>{props.v.slice(last, m.index)}</span>);
    last = m.index + m[0].length;
    text.push(<Tag key={text.length} variant="subtle" colorScheme="cyan">{m[2]}</Tag>);
  }
  text.push(<span key={text.length}>{props.v.slice(last, props.v.length)}</span>);
  return (<>
    <span style={{ whiteSpace: "pre" }}>{text}</span>
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

  advancedMode: boolean;
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
      advancedMode: ifSaved('advancedMode', false),
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
    localStorage.removeItem("NewRedirectFrom");
    console.log(this.state);
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
              <SwitchControl name="advancedMode" checked={this.state.advancedMode}
                title="Advanced mode" doc={<>
                  Toggle advanced mode.
                  <OrderedList>
                    <ListItem>Simple mode: Just send data in body as JSON object to webhook</ListItem>
                    <ListItem>Advanced mode: Template entire HTTP request, URL, Method, Headers, Body</ListItem>
                  </OrderedList>
                </>}
                on={(name, e) => this.onSwitch(name, e)} />
              <Divider mb="2rem" />
              {this.state.advancedMode ? <>
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
              <Heading as="h4" size="md">FromURI:</Heading>
              <DisplayTemplate v={this.state.FromURI} />
              <Heading as="h4" size="md">RedirectAfter:</Heading>
              <DisplayTemplate v={this.state.RedirectAfter} />
              {this.state.advancedMode ? <>
                <Heading as="h4" size="md">URLTemplate:</Heading>
                <DisplayTemplate v={this.state.URLTemplate} />
                <Heading as="h4" size="md">MethodTemplate:</Heading>
                <DisplayTemplate v={this.state.MethodTemplate} />
                <Heading as="h4" size="md">HeadersTemplate:</Heading>
                <DisplayTemplate v={this.state.HeadersTemplate} />
                <Heading as="h4" size="md">BodyTemplate:</Heading>
                <DisplayTemplate v={this.state.BodyTemplate} />
              </> : <>
                <Heading as="h4" size="md">ToURL:</Heading>
                <DisplayTemplate v={this.state.ToURL} />
              </>}
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