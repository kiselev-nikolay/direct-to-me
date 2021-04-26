import React from 'react';

import {
  Box,
  Skeleton,
  Stack,
} from '@chakra-ui/react';

import GetRedirects, { Redirect } from './API/ListRedirects';
import RedirectUI from './ui/Redirect';

let redirects: Array<Redirect> = [];

interface ListRedirectsProps { }
interface ListRedirectsState {
  redirects: Array<Redirect>;
}

export default class ListRedirects extends React.Component<ListRedirectsProps, ListRedirectsState> {
  constructor(props: ListRedirectsProps) {
    super(props);
    this.state = { redirects: [] };
  }
  componentDidMount() {
    if (redirects.length === 0) {
      setTimeout(() => {
        GetRedirects().then((rs: Array<Redirect>) => { redirects = rs; this.setState({ redirects: redirects }); });
      }, 500);
    } else {
      this.setState({ redirects: redirects });
    }
  }
  render() {
    return (<>
      {redirects.length === 0 &&
        <Stack>
          <Skeleton height="20px" />
          <Skeleton height="20px" />
          <Skeleton height="20px" />
        </Stack>
      }
      {redirects.length !== 0 &&
        <Stack>
          {redirects.map((x, i) => <Box key={i} borderWidth="1px" borderRadius="lg" p="1rem" overflow="hidden">
            <RedirectUI redirect={x} />
          </Box>)}
        </Stack>
      }
    </>);
  }
}