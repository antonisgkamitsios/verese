import { parseHTML } from './helper.ts';
import { HTMXBeforeSwapDetail } from './types';

const ASSETS_SELECTOR = `[data-track="reload"]`;

/**
 * assetWatcher listens to all htmx:beforeSwap events and checks if the current head's
 * elements tagged with `data-track="reload"` are equal with the response's head element tagged
 * with this attribute.
 * This is to make sure that we have the most current version of an asset in a case we are already on a page
 * and we click on a boosted link
 */
function assetsWatcher() {
  const stopSwapAndReload = (e: CustomEvent<HTMXBeforeSwapDetail>) => {
    // e.detail.shouldSwap = false;
    // window.location.href = e.detail.xhr.responseURL;
    e.preventDefault();
    window.location.href = e.detail.ctx.response.raw.url;
  };

  document.addEventListener('htmx:historyCacheMissLoadError', () => {
    window.location.reload();
  });

  document.addEventListener('htmx:before:swap', ((e: CustomEvent<HTMXBeforeSwapDetail>) => {
    const isBoosted = e.detail.ctx.request.headers['HX-Boosted'] === 'true';
    if (!isBoosted) {
      return;
    }

    const newHead = parseHTML(e.detail.ctx.text).querySelector('head');

    if (!newHead) {
      return;
    }

    const currentHead = document.querySelector('head')!;

    const currentSrcToAssetNodes = new Map<string, Element>();
    currentHead.querySelectorAll(ASSETS_SELECTOR).forEach((el) => {
      currentSrcToAssetNodes.set(el.outerHTML, el);
    });
    const responseTrackAssets = Array.from(newHead.querySelectorAll(ASSETS_SELECTOR));

    if (currentSrcToAssetNodes.size !== responseTrackAssets.length) {
      stopSwapAndReload(e);
      return;
    }

    for (let i = 0; i < responseTrackAssets.length; i++) {
      if (!currentSrcToAssetNodes.has(responseTrackAssets[i].outerHTML)) {
        stopSwapAndReload(e);
        return;
      }
    }
  }) as EventListener);
}
export { assetsWatcher };
