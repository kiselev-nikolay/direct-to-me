import React from 'react';

import { ArrowBackIcon } from '@chakra-ui/icons';
import {
  Box,
  Heading,
  Skeleton,
  Stack,
  Text,
} from '@chakra-ui/react';

import GetRedirects, { Redirect } from './API/ListRedirects';
import RedirectUI from './ui/Redirect';

let redirects: Array<Redirect> = [];

interface ListRedirectsProps { }
interface ListRedirectsState {
  loading: boolean;
  redirects: Array<Redirect>;
}

export default class ListRedirects extends React.Component<ListRedirectsProps, ListRedirectsState> {
  constructor(props: ListRedirectsProps) {
    super(props);
    this.state = { loading: true, redirects: [] };
  }
  componentDidMount() {
    if (redirects.length === 0) {
      setTimeout(() => {
        GetRedirects().then((rs: Array<Redirect>) => { redirects = rs; this.setState({ loading: false, redirects: redirects }); });
      }, 500);
    } else {
      this.setState({ loading: false, redirects: redirects });
    }
  }
  render() {
    return (<>
      {this.state.loading &&
        <Stack>
          <Skeleton height="40px" />
          <Skeleton height="20px" />
          <Skeleton height="20px" width="90%" />
          <Skeleton height="20px" />
          <Skeleton height="20px" width="85%" />
          <Skeleton height="20px" width="95%" />
          <Skeleton height="20px" />
        </Stack>
      }
      {!this.state.loading && redirects.length === 0 && <Box m="2rem">
        <Heading mb="1rem">Empty yet!</Heading>
        <Text><ArrowBackIcon /> Try creating a new redirect in the "New redirect" section of the left navigation.</Text>
      </Box>}
      {!this.state.loading && redirects.length !== 0 &&
        <Stack>
          {redirects.map((x, i) => <Box key={i} borderWidth="1px" borderRadius="lg" p="1rem" overflow="hidden">
            <RedirectUI redirect={x} />
          </Box>)}
        </Stack>
      }
    </>);
  }
}