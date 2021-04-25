import React from 'react';

import {
  Heading,
  Tag,
} from '@chakra-ui/react';

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
interface RedirectProps {
  redirect: {
    FromURI: string;
    RedirectAfter: string;

    isAdvanced: boolean;

    URLTemplate?: string;
    MethodTemplate?: string;
    HeadersTemplate?: string;
    BodyTemplate?: string;

    ToURL?: string;
  };
}
export default function Redirect(props: RedirectProps) {
  return (<>
    <Heading as="h4" size="md">FromURI:</Heading>
    <DisplayTemplate v={props.redirect.FromURI} />
    <Heading as="h4" size="md">RedirectAfter:</Heading>
    <DisplayTemplate v={props.redirect.RedirectAfter} />
    {props.redirect.isAdvanced ? <>
      <Heading as="h4" size="md">URLTemplate:</Heading>
      <DisplayTemplate v={props.redirect.URLTemplate} />
      <Heading as="h4" size="md">MethodTemplate:</Heading>
      <DisplayTemplate v={props.redirect.MethodTemplate} />
      <Heading as="h4" size="md">HeadersTemplate:</Heading>
      <DisplayTemplate v={props.redirect.HeadersTemplate} />
      <Heading as="h4" size="md">BodyTemplate:</Heading>
      <DisplayTemplate v={props.redirect.BodyTemplate} />
    </> : <>
      <Heading as="h4" size="md">ToURL:</Heading>
      <DisplayTemplate v={props.redirect.ToURL} />
    </>}
  </>);
}