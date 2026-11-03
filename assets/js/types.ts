export type HTMXBeforeSwapDetail = {
  ctx: {
    request: {
      headers: { [key: string]: string };
    };
    response: {
      raw: Response;
    };
    swap: string;
    text: string;
  };
};
