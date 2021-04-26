import React from 'react';

import {
  Table,
  Tag,
  Tbody,
  Td,
  Tr,
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
    text.push(<Tag key={text.length} variant="subtle" colorScheme="cyan" mt={.6}>{m[2]}</Tag>);
  }
  text.push(<span key={text.length}>{props.v.slice(last, props.v.length)}</span>);
  return (<>
    <span style={{ whiteSpace: "pre" }}>{text}</span>
  </>);
}
interface StatLineProps {
  title: string;
  value: string;
}
function StatLine(props: StatLineProps) {
  return (<Tr>
    <Td width="20%">{props.title}</Td>
    <Td width="80%" lineHeight="1.8"><DisplayTemplate v={props.value} /></Td>
  </Tr>);
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
    <Table variant="simple">
      <Tbody>
        <StatLine title="From URI" value={props.redirect.FromURI} />
        <StatLine title="Redirect after" value={props.redirect.RedirectAfter} />
        {props.redirect.isAdvanced ? <>
          <StatLine title="URL template" value={props.redirect.URLTemplate} />
          <StatLine title="Method template" value={props.redirect.MethodTemplate} />
          <StatLine title="Headers template" value={props.redirect.HeadersTemplate} />
          <StatLine title="Body template" value={props.redirect.BodyTemplate} />
        </> : <>
          <StatLine title="Send JSON to URL" value={props.redirect.ToURL} />
        </>}
      </Tbody>
    </Table>
  </>);
}