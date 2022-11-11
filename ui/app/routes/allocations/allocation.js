import Route from '@ember/routing/route';
import { inject as service } from '@ember/service';
import { collect } from '@ember/object/computed';
import {
  watchRecord,
  watchNonStoreRecords,
} from 'nomad-ui/utils/properties/watch';
import WithWatchers from 'nomad-ui/mixins/with-watchers';
import notifyError from 'nomad-ui/utils/notify-error';
export default class AllocationRoute extends Route.extend(WithWatchers) {
  @service store;

  startWatchers(controller, model) {
    if (model) {
      controller.set('watcher', this.watch.perform(model));

      const anyGroupServicesAreNomad = !!model.taskGroup?.services?.filterBy(
        'provider',
        'nomad'
      ).length;

      const anyTaskServicesAreNomad = model.states
        .mapBy('task.services')
        .compact()
        .map((fragmentClass) => fragmentClass.mapBy('provider'))
        .flat()
        .any((provider) => provider === 'nomad');

      // Conditionally Long Poll /checks endpoint if alloc has nomad services
      if (anyGroupServicesAreNomad || anyTaskServicesAreNomad) {
        controller.set(
          'watchHealthChecks',
          this.watchHealthChecks.perform(model, 'getServiceHealth', 2000)
        );
      }
    }
  }

  async model() {
    try {
      // Preload the job for the allocation since it's required for the breadcrumb trail
      const allocation = await super.model(...arguments);
      const jobId = allocation?.belongsTo('job').id();
      const getJob = this.store.findRecord('job', jobId);
      const getNamespaces = this.store.findAll('namespace');
      await Promise.all([getJob, getNamespaces]);
      return allocation;
    } catch (e) {
      notifyError.call(this, e);
    }
  }

  @watchRecord('allocation') watch;
  @watchNonStoreRecords('allocation') watchHealthChecks;

  @collect('watch', 'watchHealthChecks') watchers;
}
