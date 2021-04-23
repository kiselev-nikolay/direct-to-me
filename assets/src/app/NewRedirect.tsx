import React from 'react';

import { QuestionIcon } from '@chakra-ui/icons';
import {
  Box,
  Button,
  Flex,
  FormControl,
  FormLabel,
  Input,
  Popover,
  PopoverArrow,
  PopoverBody,
  PopoverCloseButton,
  PopoverContent,
  PopoverTrigger,
  Switch,
  Textarea,
} from '@chakra-ui/react';

interface LabelWithHelpProps {
  title: string;
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
        <PopoverCloseButton />
        <PopoverBody>{props.children}</PopoverBody>
      </PopoverContent>
    </Popover>
  );
}

interface ControlProps {
  title: string;
  doc: string;
  name: string;
  on: (name: string, e: any) => void;
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
  AdvancedMode: boolean;
}
class NewRedirectFrom extends React.Component<NewRedirectFromProps, NewRedirectFromState> {
  constructor(props: NewRedirectFromProps) {
    super(props);
    this.state = {
      FromURI: localStorage.getItem('FromURI') || '',
      ToURL: localStorage.getItem('ToURL') || '',
      RedirectAfter: localStorage.getItem('RedirectAfter') || '',
      URLTemplate: localStorage.getItem('URLTemplate') || '',
      MethodTemplate: localStorage.getItem('MethodTemplate') || '',
      HeadersTemplate: localStorage.getItem('HeadersTemplate') || '',
      BodyTemplate: localStorage.getItem('BodyTemplate') || '',
      AdvancedMode: localStorage.getItem('AdvancedMode') === 'true' || false,
    };
  }
  onChange(name: string, e: any) {
    let newState: any = {};
    newState[name] = e.target.value;
    this.setState(newState);
    e.preventDefault();
  }
  onSwitch(name: string, e: any) {
    let newState: any = {};
    newState[name] = e.target.checked;
    this.setState(newState);
    e.preventDefault();
  }
  onSubmit(e: any) {
    console.log(this.state);
    e.preventDefault();
  }
  render() {
    function LineControl(props: ControlProps) {
      return (<>
        <FormControl mb="1rem">
          <LabelWithHelp title={props.title}>{props.doc}</LabelWithHelp>
          <Input name={props.name} onChange={e => props.on(props.name, e)} />
        </FormControl>
      </>);
    }
    function AreaControl(props: ControlProps) {
      return (<>
        <FormControl mb="1rem">
          <LabelWithHelp title={props.title}>{props.doc}</LabelWithHelp>
          <Textarea name={props.name} onChange={e => props.on(props.name, e)} />
        </FormControl>
      </>);
    }
    function SwitchControl(props: ControlProps) {
      return (<>
        <FormControl display="flex" alignItems="center" mb="1rem">
          <LabelWithHelp title={props.title}>{props.doc}</LabelWithHelp>
          <Switch onChange={e => props.on(props.name, e)} />
        </FormControl>
      </>);
    }
    return (
      <>
        <Flex>
          <Box flex="1">
            <Box w="100%" pr="1rem">
              <LineControl name="FromURI" title="From URI" doc=""
                on={(name, e) => this.onChange(name, e)} />
              <LineControl name="RedirectAfter" title="Redirect after" doc=""
                on={(name, e) => this.onChange(name, e)} />
            </Box>
          </Box>
          <Box flex="1">
            <Box w="100%" pl="1rem">
              {this.state.AdvancedMode ? <>
                <LineControl name="URLTemplate" title="URL template" doc=""
                  on={(name, e) => this.onChange(name, e)} />
                <LineControl name="MethodTemplate" title="Method template" doc=""
                  on={(name, e) => this.onChange(name, e)} />
                <AreaControl name="HeadersTemplate" title="Headers template" doc=""
                  on={(name, e) => this.onChange(name, e)} />
                <AreaControl name="BodyTemplate" title="Body template" doc=""
                  on={(name, e) => this.onChange(name, e)} />
              </> : <>
                <LineControl name="ToURL" title="Send JSON to URL in background" doc=""
                  on={(name, e) => this.onChange(name, e)} />
              </>
              }
            </Box>
          </Box>
        </Flex>
        <SwitchControl name="AdvancedMode" title="Advanced mode" doc=""
          on={(name, e) => this.onSwitch(name, e)} />
        <Button colorScheme="green" onClick={e => this.onSubmit(e)}>Submit</Button>
      </>
    );
  }
}

export default NewRedirectFrom;